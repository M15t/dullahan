package session

import (
	"dullahan/internal/model"
	"fmt"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
)

func (s *Session) getTotalIncome(session *model.Session) float64 {
	var totalIncome float64 = 0.0
	if len(session.Incomes) > 0 {
		s.db.Income.SumTotalIncome(s.db.GDB, &totalIncome, session.ID)
	}

	return totalIncome
}

func (s *Session) getTotalMonthlyPaymentDebt(session *model.Session) float64 {
	var totalMonthlyPaymentDebt float64 = 0.0
	if len(session.Debts) > 0 {
		s.db.Debt.SumTotalMonthlyPaymentDebt(s.db.GDB, &totalMonthlyPaymentDebt, session.ID)
	}

	return totalMonthlyPaymentDebt
}

func (s *Session) getTotalRemaingingDebt(session *model.Session) float64 {
	var totalRemaingingDebt float64 = 0.0
	if len(session.Debts) > 0 {
		s.db.Debt.SumTotalRemainingDebt(s.db.GDB, &totalRemaingingDebt, session.ID)
	}

	return totalRemaingingDebt
}

func (s *Session) calculateSession(cache *bigcache.BigCache, session *model.Session) error {
	// * init first node
	node := calculateNode(session.ID, "0",
		roundFloat(session.CurrentBalance),
		roundFloat(session.TotalEssentialExpense),
		roundFloat(session.TotalNonEssentialExpense),
		roundFloat(s.getTotalRemaingingDebt(session)),
		roundFloat(s.getTotalIncome(session)),
		roundFloat(s.getTotalMonthlyPaymentDebt(session)))

	// * set node to cache
	setNodeToCache(cache, node)

	// * return latest information
	session.TotalAllIncome = node.TotalAllIncome
	session.TotalAllExpense = node.TotalAllExpense
	session.TotalMonthlyPaymentDebt = node.TotalMonthlyPaymentDebt
	session.MonthlyNetFlow = node.MonthlyNetFlow
	session.ExpectedEmergencyFund = node.ExpectedEmergencyFund
	session.ExpectedRainydayFund = node.ExpectedRainydayFund
	session.ActualEmergencyFund = node.ActualEmergencyFund
	session.ActualRainydayFund = node.ActualRainydayFund
	session.FunFund = node.FunFund
	session.RetirementPlan = node.RetirementPlan
	session.IsAchivedEmergencyFund = node.IsAchivedEmergencyFund
	session.IsAchivedRainydayFund = node.IsAchivedRainydayFund
	session.IsAchivedInvestment = node.IsAchivedInvestment
	session.IsAchivedRetirementPlan = node.IsAchivedRetirementPlan
	session.Status = node.Status
	session.Description = node.Descrtiption

	// * update session
	return s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"total_all_income":           node.TotalAllIncome,
		"total_all_expense":          node.TotalAllExpense,
		"total_monthly_payment_debt": node.TotalMonthlyPaymentDebt,
		"monthly_net_flow":           node.MonthlyNetFlow,
		"status":                     node.Status,
		"expected_emergency_fund":    node.ExpectedEmergencyFund,
		"expected_rainyday_fund":     node.ExpectedRainydayFund,
		"actual_emergency_fund":      node.ActualEmergencyFund,
		"actual_rainyday_fund":       node.ActualRainydayFund,
		"fun_fund":                   node.FunFund,
		"retirement_plan":            node.RetirementPlan,
		"is_achived_emergency_fund":  node.IsAchivedEmergencyFund,
		"is_achived_rainyday_fund":   node.IsAchivedRainydayFund,
		"is_achived_investment":      node.IsAchivedInvestment,
		"is_achived_retirement_plan": node.IsAchivedRetirementPlan,
	}, session.ID)
}

func calculateDebtPaidEachMonth(monthlyPayment, annualInterest float64) float64 {
	interestPaid := roundFloat((annualInterest / 12.0) * monthlyPayment)
	return roundFloat(monthlyPayment-interestPaid) * 3
}

func calculateNode(sessionID int64,
	nodeName string,
	currentAsset,
	totalEssentialExpense,
	totalNonEssentialExpense,
	totalRemainingDebt,
	totalIncome,
	totalMonthlyPaymentDebt float64) *model.DataNode {
	var funFund float64
	var isAchivedInvestment bool

	totalExpense := roundFloat(totalEssentialExpense + totalNonEssentialExpense)

	monthlyNetFlow := totalIncome - (totalMonthlyPaymentDebt + totalExpense)

	// * total asset this month <-
	currentAsset = currentAsset + monthlyNetFlow

	netAsset := currentAsset - totalRemainingDebt

	expectedEmergencyFund := roundFloat(totalEssentialExpense * EmergencyFundRate)
	expectedRainydayFund := roundFloat(totalEssentialExpense * RainydayFundRate)

	actualEmergencyFund := roundFloat(netAsset)
	if actualEmergencyFund <= 0 {
		actualEmergencyFund = 0
	}
	actualRainydayFund := roundFloat(netAsset - expectedEmergencyFund)
	if actualRainydayFund <= 0 {
		actualRainydayFund = 0
	}

	retirementPlan := roundFloat(totalEssentialExpense * 12 * RetirementPlanRate)

	isAchivedEmergencyFund := netAsset >= expectedEmergencyFund
	if isAchivedEmergencyFund {
		actualEmergencyFund = expectedEmergencyFund
	}
	isAchivedRainydayFund := netAsset >= (expectedEmergencyFund + expectedRainydayFund)
	if isAchivedRainydayFund {
		actualRainydayFund = expectedRainydayFund
	}
	isAchivedRetirementPlan := netAsset >= retirementPlan && retirementPlan > 0

	if isAchivedEmergencyFund && isAchivedRainydayFund {
		funFund = monthlyNetFlow / 100 * 20
		isAchivedInvestment = true

		// * calculate R
		r := currentAsset - (actualEmergencyFund + actualRainydayFund)

		if r >= 0 {
			currentAsset = currentAsset + (r * UserBudgetIncreasementRate)
		}
	}

	var status string
	switch {
	case monthlyNetFlow < 0:
		status = model.SessionStatusBD
	case 0 <= monthlyNetFlow && monthlyNetFlow < BankruptCeil:
		status = model.SessionStatusPC2PC
	case BankruptCeil <= monthlyNetFlow && monthlyNetFlow <= totalEssentialExpense:
		status = model.SessionStatusLFF
	case totalEssentialExpense < monthlyNetFlow:
		status = model.SessionStatusGFF
	default:
		status = model.SessionStatusDefault
	}

	return &model.DataNode{
		SessionID:               sessionID,
		NodeName:                nodeName,
		CurrentAsset:            roundFloat(currentAsset),
		TotalAllIncome:          totalIncome,
		TotalAllExpense:         totalExpense,
		TotalMonthlyPaymentDebt: roundFloat(totalMonthlyPaymentDebt),
		MonthlyNetFlow:          monthlyNetFlow,
		Status:                  status,
		Descrtiption:            model.SessionStatusDescriptions[status],
		ExpectedEmergencyFund:   expectedEmergencyFund,
		ExpectedRainydayFund:    expectedRainydayFund,
		ActualEmergencyFund:     actualEmergencyFund,
		ActualRainydayFund:      actualRainydayFund,
		FunFund:                 funFund,
		RetirementPlan:          retirementPlan,
		IsAchivedEmergencyFund:  isAchivedEmergencyFund,
		IsAchivedRainydayFund:   isAchivedRainydayFund,
		IsAchivedInvestment:     isAchivedInvestment,
		IsAchivedRetirementPlan: isAchivedRetirementPlan,
	}
}

func setNodeToCache(cache *bigcache.BigCache, node *model.DataNode) {
	cache.Set(fmt.Sprintf("%d_%s_current_asset", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.CurrentAsset)))
	cache.Set(fmt.Sprintf("%d_%s_total_all_income", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.TotalAllIncome)))
	cache.Set(fmt.Sprintf("%d_%s_total_all_expense", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.TotalAllExpense)))
	cache.Set(fmt.Sprintf("%d_%s_total_monthly_payment_debt", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.TotalMonthlyPaymentDebt)))
	cache.Set(fmt.Sprintf("%d_%s_monthly_net_flow", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.MonthlyNetFlow)))

	cache.Set(fmt.Sprintf("%d_%s_status", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%s", node.Status)))
	cache.Set(fmt.Sprintf("%d_%s_description", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%s", node.Descrtiption)))

	cache.Set(fmt.Sprintf("%d_%s_expected_emergency_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.ExpectedEmergencyFund)))
	cache.Set(fmt.Sprintf("%d_%s_actual_emergency_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.ActualEmergencyFund)))

	cache.Set(fmt.Sprintf("%d_%s_expected_rainyday_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.ExpectedRainydayFund)))
	cache.Set(fmt.Sprintf("%d_%s_actual_rainyday_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.ActualRainydayFund)))

	cache.Set(fmt.Sprintf("%d_%s_fun_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.FunFund)))
	cache.Set(fmt.Sprintf("%d_%s_retirement_plan", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.RetirementPlan)))

	cache.Set(fmt.Sprintf("%d_%s_is_achived_emergency_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedEmergencyFund)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_rainyday_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedRainydayFund)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_investment", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedInvestment)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_retirement_plan", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedRetirementPlan)))
}

func getNodeFromCache(cache *bigcache.BigCache, session *model.Session, nodeName string) *model.DataNode {
	cca, _ := cache.Get(fmt.Sprintf("%d_%s_current_asset", session.ID, nodeName))
	cCurrentAsset, _ := strconv.ParseFloat(string(cca), 64)

	ctai, _ := cache.Get(fmt.Sprintf("%d_%s_total_all_income", session.ID, nodeName))
	cTotalAllIncome, _ := strconv.ParseFloat(string(ctai), 64)

	ctae, _ := cache.Get(fmt.Sprintf("%d_%s_total_all_expense", session.ID, nodeName))
	cTotalAllExpense, _ := strconv.ParseFloat(string(ctae), 64)

	ctmpd, _ := cache.Get(fmt.Sprintf("%d_%s_total_monthly_payment_debt", session.ID, nodeName))
	cTotalMonthlyPaymentDebt, _ := strconv.ParseFloat(string(ctmpd), 64)

	cmnf, _ := cache.Get(fmt.Sprintf("%d_%s_monthly_net_flow", session.ID, nodeName))
	cMonthlyNetFlow, _ := strconv.ParseFloat(string(cmnf), 64)

	cStatus, _ := cache.Get(fmt.Sprintf("%d_%s_status", session.ID, nodeName))

	cDescription, _ := cache.Get(fmt.Sprintf("%d_%s_description", session.ID, nodeName))

	ceef, _ := cache.Get(fmt.Sprintf("%d_%s_expected_emergency_fund", session.ID, nodeName))
	cExpectedEmergencyFund, _ := strconv.ParseFloat(string(ceef), 64)

	cerf, _ := cache.Get(fmt.Sprintf("%d_%s_expected_rainyday_fund", session.ID, nodeName))
	cExpectedRainydayFund, _ := strconv.ParseFloat(string(cerf), 64)

	caef, _ := cache.Get(fmt.Sprintf("%d_%s_actual_emergency_fund", session.ID, nodeName))
	cActualEmergencyFund, _ := strconv.ParseFloat(string(caef), 64)

	carf, _ := cache.Get(fmt.Sprintf("%d_%s_actual_rainyday_fund", session.ID, nodeName))
	cActualRainydayFund, _ := strconv.ParseFloat(string(carf), 64)

	cff, _ := cache.Get(fmt.Sprintf("%d_%s_fun_fund", session.ID, nodeName))
	cFunFund, _ := strconv.ParseFloat(string(cff), 64)

	crp, _ := cache.Get(fmt.Sprintf("%d_%s_retirement_plan", session.ID, nodeName))
	cRetirementPlan, _ := strconv.ParseFloat(string(crp), 64)

	ciae, _ := cache.Get(fmt.Sprintf("%d_%s_is_achived_emergency_fund", session.ID, nodeName))
	cIsAchivedEmergencyFund, _ := strconv.ParseBool(string(ciae))

	ciarf, _ := cache.Get(fmt.Sprintf("%d_%s_is_achived_rainyday_fund", session.ID, nodeName))
	cIsAchivedRainydayFund, _ := strconv.ParseBool(string(ciarf))

	ciai, _ := cache.Get(fmt.Sprintf("%d_%s_is_achived_investment", session.ID, nodeName))
	cIsAchivedInvestment, _ := strconv.ParseBool(string(ciai))

	ciarp, _ := cache.Get(fmt.Sprintf("%d_%s_is_achived_retirement_plan", session.ID, nodeName))
	cIsAchivedRetirementPlan, _ := strconv.ParseBool(string(ciarp))

	return &model.DataNode{
		SessionID:               session.ID,
		NodeName:                nodeName,
		CurrentAsset:            cCurrentAsset,
		TotalAllIncome:          cTotalAllIncome,
		TotalAllExpense:         cTotalAllExpense,
		TotalMonthlyPaymentDebt: cTotalMonthlyPaymentDebt,
		MonthlyNetFlow:          cMonthlyNetFlow,
		Status:                  string(cStatus),
		Descrtiption:            string(cDescription),

		ExpectedEmergencyFund: cExpectedEmergencyFund,
		ExpectedRainydayFund:  cExpectedRainydayFund,
		ActualEmergencyFund:   cActualEmergencyFund,
		ActualRainydayFund:    cActualRainydayFund,

		FunFund:        cFunFund,
		RetirementPlan: cRetirementPlan,

		IsAchivedEmergencyFund:  cIsAchivedEmergencyFund,
		IsAchivedRainydayFund:   cIsAchivedRainydayFund,
		IsAchivedInvestment:     cIsAchivedInvestment,
		IsAchivedRetirementPlan: cIsAchivedRetirementPlan,
	}
}

func setDebtNodeToCache(cache *bigcache.BigCache, node *model.DataDebtNode) {
	cache.Set(fmt.Sprintf("%d_%d_%d_%s_remaining_amount", node.SessionID, node.DebtID, node.Index, node.NodeName), []byte(fmt.Sprintf("%f", node.RemainingAmount)))
	cache.Set(fmt.Sprintf("%d_%d_%d_%s_monthly_payment", node.SessionID, node.DebtID, node.Index, node.NodeName), []byte(fmt.Sprintf("%f", node.MonthlyPayment)))
	cache.Set(fmt.Sprintf("%d_%d_%d_%s_is_eligible_paid_off", node.SessionID, node.DebtID, node.Index, node.NodeName), []byte(fmt.Sprintf("%v", node.IsEligiblePaidOff)))
	cache.Set(fmt.Sprintf("%d_%d_%d_%s_is_paid_off", node.SessionID, node.DebtID, node.Index, node.NodeName), []byte(fmt.Sprintf("%v", node.IsPaidOff)))
}

func getDebtNodeFromCache(cache *bigcache.BigCache, session *model.Session, debt *model.Debt, index int, nodeName string) *model.DataDebtNode {
	ra, _ := cache.Get(fmt.Sprintf("%d_%d_%d_%s_remaining_amount", session.ID, debt.ID, index, nodeName))
	remainingAmount, _ := strconv.ParseFloat(string(ra), 64)

	mp, _ := cache.Get(fmt.Sprintf("%d_%d_%d_%s_monthly_payment", session.ID, debt.ID, index, nodeName))
	monthlyPayment, _ := strconv.ParseFloat(string(mp), 64)

	iepo, _ := cache.Get(fmt.Sprintf("%d_%d_%d_%s_is_eligible_paid_off", session.ID, debt.ID, index, nodeName))
	isEligiblePaidOff, _ := strconv.ParseBool(string(iepo))

	ipo, _ := cache.Get(fmt.Sprintf("%d_%d_%d_%s_is_paid_off", session.ID, debt.ID, index, nodeName))
	isPaidOff, _ := strconv.ParseBool(string(ipo))

	return &model.DataDebtNode{
		SessionID: session.ID,
		NodeName:  nodeName,
		DebtID:    debt.ID,
		Index:     index,

		RemainingAmount:   remainingAmount,
		MonthlyPayment:    monthlyPayment,
		IsEligiblePaidOff: isEligiblePaidOff,
		IsPaidOff:         isPaidOff,
	}
}

func generateDatasets(cache *bigcache.BigCache, rec *model.Session) []*model.DataSet {
	startDate := time.Now()
	endDate := time.Date(CustomYear, CustomMonth, CustomDay, 0, 0, 0, 0, time.UTC)

	if len(rec.Debts) > 0 { // * case with debt
		return generateDataSetWithDebt(cache, rec, startDate, endDate)
	}

	// * case without debt
	return generateDataSetWithoutDebt(cache, rec, startDate, endDate)
}

func generateDataSetWithoutDebt(cache *bigcache.BigCache, rec *model.Session, startDate, endDate time.Time) []*model.DataSet {
	datasets := make([]*model.DataSet, 0)

	for i, q := range generateMonths(startDate, endDate) {
		fmt.Println(q)
		curNode := new(model.DataNode)
		if i == 0 {
			curNode = getNodeFromCache(cache, rec, "0")
		} else {
			prevNode := getNodeFromCache(cache, rec, fmt.Sprintf("%d", i-1))

			curNode = calculateNode(rec.ID, fmt.Sprintf("%d", i),
				roundFloat(prevNode.CurrentAsset),        // * dynamic
				roundFloat(rec.TotalEssentialExpense),    // * static
				roundFloat(rec.TotalNonEssentialExpense), // * static
				0.0,                                      // * no debt
				roundFloat(rec.TotalAllIncome),           // * static
				0.0)                                      // * no debt

			setNodeToCache(cache, curNode)
		}

		if curNode.CurrentAsset > 0 {
			datasets = append(datasets, &model.DataSet{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: curNode.CurrentAsset,
			})
		} else {
			datasets = append(datasets, &model.DataSet{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: 0,
			})
			break
		}
	}

	return datasets
}

func generateDataSetWithDebt(cache *bigcache.BigCache, rec *model.Session, startDate, endDate time.Time) []*model.DataSet {
	datasets := make([]*model.DataSet, 0)
	eligiblePaidOff := make(map[int]bool)
	eligiblePaidOff[0] = true

	for i := 1; i <= len(rec.Debts); i++ {
		eligiblePaidOff[i] = false
	}

	for i, q := range generateMonths(startDate, endDate) {
		var curNode *model.DataNode
		if i == 0 {
			curNode = getNodeFromCache(cache, rec, "0")
		}

		prevNode := getNodeFromCache(cache, rec, fmt.Sprintf("%d", i-1))

		for j, debt := range rec.Debts {
			var currentRemainingDebt float64
			var isPaidOff bool
			if i == 0 {
				currentRemainingDebt = debt.RemainingAmount

				setDebtNodeToCache(cache, &model.DataDebtNode{
					SessionID:       rec.ID,
					NodeName:        "0",
					DebtID:          debt.ID,
					Index:           j,
					RemainingAmount: currentRemainingDebt,
					MonthlyPayment:  debt.MonthlyPayment,
					IsPaidOff:       isPaidOff,
				})
			} else {
				prevDebtNode := getDebtNodeFromCache(cache, rec, debt, j, fmt.Sprintf("%d", i-1))

				currentRemainingDebt = roundFloat(prevDebtNode.RemainingAmount - (calculateDebtPaidEachMonth(debt.MonthlyPayment, debt.AnnualInterest)))

				// * paid off
				if currentRemainingDebt > 0 && prevNode.CurrentAsset-currentRemainingDebt >= 0 && eligiblePaidOff[j] {
					prevNode.CurrentAsset = prevNode.CurrentAsset - currentRemainingDebt
					currentRemainingDebt = 0
					isPaidOff = true

					// * update next debt to be eligible paid off
					eligiblePaidOff[j+1] = true
				}

				setDebtNodeToCache(cache, &model.DataDebtNode{
					SessionID:       rec.ID,
					NodeName:        fmt.Sprintf("%d", i),
					DebtID:          debt.ID,
					Index:           j,
					RemainingAmount: currentRemainingDebt,
					MonthlyPayment:  debt.MonthlyPayment,
					IsPaidOff:       isPaidOff,
				})

				if prevDebtNode.RemainingAmount <= 0 {
					prevNode.CurrentAsset = prevNode.CurrentAsset + debt.MonthlyPayment
				}

				curNode = calculateNode(rec.ID, fmt.Sprintf("%d", i),
					roundFloat(prevNode.CurrentAsset),        // * dynamic
					roundFloat(rec.TotalEssentialExpense),    // * static
					roundFloat(rec.TotalNonEssentialExpense), // * static
					currentRemainingDebt,                     // * dynamic
					roundFloat(rec.TotalAllIncome),           // *static
					debt.MonthlyPayment)                      // * dynamic
			}

			// * append debt
			if currentRemainingDebt >= 0 {
				datasets = append(datasets, &model.DataSet{
					Group: debt.Name,
					Key:   getMonth(startDate, q),
					Debt:  roundFloat(currentRemainingDebt),
				})
			}
		}

		setNodeToCache(cache, curNode)

		// * append asset
		if curNode.CurrentAsset > 0 {
			datasets = append(datasets, &model.DataSet{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: curNode.CurrentAsset,
			})
		} else {
			datasets = append(datasets, &model.DataSet{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: 0,
			})
			break
		}
	}

	return datasets
}
