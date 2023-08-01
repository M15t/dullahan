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

func (s *Session) getTotalExpense(session *model.Session) float64 {
	var totalExpense float64 = 0.0
	if len(session.Expenses) > 0 {
		s.db.Expense.SumTotalExpense(s.db.GDB, &totalExpense, session.ID)
	}

	return totalExpense
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
	var isPaidAllDebt bool

	if len(session.Debts) == 0 {
		isPaidAllDebt = true
	}

	session.TotalAllIncome = s.getTotalIncome(session)
	session.TotalMonthlyPaymentDebt = s.getTotalMonthlyPaymentDebt(session)
	session.TotalAllExpense = s.getTotalExpense(session)

	// * init first node
	node := calculateNode(session, 0, "0",
		roundFloat(session.CurrentBalance), 0, isPaidAllDebt)

	// * set node to cache
	setNodeToCache(cache, node)

	// * return latest information
	session.TotalAllIncome = node.TotalAllIncome
	session.TotalAllExpense = node.TotalAllExpense
	session.TotalMonthlyPaymentDebt = node.TotalMonthlyPaymentDebt
	session.MonthlyNetFlow = node.MonthlyNetFlow
	session.ExpectedEmergencyFund = node.ExpectedEmergencyFund
	session.ExpectedRainydayFund = node.ExpectedRainydayFund
	session.ExpectedFunFund = node.ExpectFunFund
	session.ActualEmergencyFund = node.ActualEmergencyFund
	session.ActualRainydayFund = node.ActualRainydayFund
	session.ActualFunFund = node.ActualFunFund
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
		"expected_fun_fund":          node.ExpectFunFund,
		"actual_emergency_fund":      node.ActualEmergencyFund,
		"actual_rainyday_fund":       node.ActualRainydayFund,
		"actual_fun_fund":            node.ActualFunFund,
		"retirement_plan":            node.RetirementPlan,
		"is_achived_emergency_fund":  node.IsAchivedEmergencyFund,
		"is_achived_rainyday_fund":   node.IsAchivedRainydayFund,
		"is_achived_investment":      node.IsAchivedInvestment,
		"is_achived_retirement_plan": node.IsAchivedRetirementPlan,
	}, session.ID)
}

func (s *Session) generateLinechart(cache *bigcache.BigCache, rec *model.Session) []*model.LineChart {
	now := time.Now()

	startDate := now
	endDate := time.Date(now.Year()+YearsForCalculation, CustomMonth, CustomDay, 0, 0, 0, 0, time.UTC)

	if len(rec.Debts) > 0 { // * case with debt
		return s.generateDatasetWithDebt(cache, rec, startDate, endDate)
	}

	// * case without debt
	return s.generateDatasetWithoutDebt(cache, rec, startDate, endDate)
}

func (s *Session) generateDatasetWithoutDebt(cache *bigcache.BigCache, rec *model.Session, startDate, endDate time.Time) []*model.LineChart {
	var monthlyNetFlowWithoutDebt, currentAssetToMillionaire float64
	var millionaireDate string
	datasets := make([]*model.LineChart, 0)
	events := make(map[int]string)
	// 0 emergency fund
	// 1 rainy day fund
	// 2 investment
	// 3 financial freedom

	for i, q := range generateMonths(startDate, endDate) {
		var curNode, prevNode *model.DataNode
		var currentAsset float64
		if i == 0 {
			currentAsset = rec.CurrentBalance + calculateMonthlyNetFlow(rec)
		} else {
			prevNode = getNodeFromCache(cache, rec, fmt.Sprintf("%d", i-1))
			currentAsset = prevNode.CurrentAsset + calculateMonthlyNetFlow(rec)
		}

		// * calculate current node
		curNode = calculateNode(rec, q, fmt.Sprintf("%d", i),
			currentAsset, // * dynamic
			0,            // * no debt
			true)         // * no debt

		if rec.TotalAllIncome > 0 {
			if prevNode != nil && !prevNode.IsAchivedEmergencyFund && curNode.IsAchivedEmergencyFund || prevNode == nil && curNode.IsAchivedEmergencyFund {
				events[0] = getMonthAndYear(startDate, q)
			}

			if prevNode != nil && !prevNode.IsAchivedRainydayFund && curNode.IsAchivedRainydayFund || prevNode == nil && curNode.IsAchivedRainydayFund {
				events[1] = getMonthAndYear(startDate, q)
			}

			if prevNode != nil && !prevNode.IsAchivedInvestment && curNode.IsAchivedEmergencyFund && curNode.IsAchivedRainydayFund || prevNode == nil && curNode.IsAchivedEmergencyFund && curNode.IsAchivedRainydayFund {
				events[2] = getMonthAndYear(startDate, q)
				monthlyNetFlowWithoutDebt = curNode.MonthlyNetFlow
				currentAssetToMillionaire = currentAsset
			}

			if prevNode != nil && !prevNode.IsAchivedFinancialFreedom && curNode.IsAchivedFinancialFreedom {
				events[3] = getMonthAndYear(startDate, q)
			}
		}

		setNodeToCache(cache, curNode)

		// * append asset
		if curNode.CurrentAsset > 0 {
			datasets = append(datasets, &model.LineChart{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: curNode.CurrentAsset,
			})
		} else {
			datasets = append(datasets, &model.LineChart{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: 0,
			})
			events[4] = getMonthAndYear(startDate, q)
			break
		}
	}

	if currentAssetToMillionaire > 0 {
		remainingAssetToMillionaire := MillionaireRate - currentAssetToMillionaire

		becomeMillionaireIn := remainingAssetToMillionaire / monthlyNetFlowWithoutDebt // months

		// fmt.Printf("becomeMillionaireIn ==== %f years \n", becomeMillionaireIn/12/2)

		t := startDate.AddDate(0, int(becomeMillionaireIn/2), 0)
		millionaireDate = t.Format("Jan 2006")
	}

	if err := s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"forecast_emergency_budget_filled_date": events[0],
		"forecast_start_investing_date":         events[2],
		"forecast_rainyday_budget_filled_date":  events[1],
		"forecast_financial_freedom_date":       events[3],
		"forecast_millionaire_date":             millionaireDate,
		"forecast_bankrupt":                     events[4],
	}, rec.ID); err != nil {
		fmt.Println("Error updating forecast events")
	}

	return datasets
}

func (s *Session) generateDatasetWithDebt(cache *bigcache.BigCache, rec *model.Session, startDate, endDate time.Time) []*model.LineChart {
	var monthlyNetFlowWithoutDebt, currentAssetToMillionaire float64
	var millionaireDate string
	datasets := make([]*model.LineChart, 0)
	eligiblePaidOff := map[int]bool{0: true}
	events := make(map[int]string)
	// 0 emergency fund
	// 1 rainy day fund
	// 2 investment
	// 3 financial freedom
	// 4 bankrupt

	for i := 1; i <= len(rec.Debts); i++ {
		eligiblePaidOff[i] = false
	}

	for i, q := range generateMonths(startDate, endDate) {
		var curNode, prevNode *model.DataNode
		var currentAsset, totalRemainingDebt float64
		if i == 0 {
			currentAsset = rec.CurrentBalance + calculateMonthlyNetFlow(rec)
		} else {
			prevNode = getNodeFromCache(cache, rec, fmt.Sprintf("%d", i-1))
			currentAsset = prevNode.CurrentAsset + calculateMonthlyNetFlow(rec)
		}

		for j, debt := range rec.Debts {
			var currentRemainingDebt, totalRemainingAmount float64
			var isPaidOff bool
			var prevDebtNode *model.DataDebtNode

			if i == 0 {
				totalRemainingAmount = debt.RemainingAmount
			} else {
				prevDebtNode = getDebtNodeFromCache(cache, rec, debt, j, fmt.Sprintf("%d", i-1))
				totalRemainingAmount = prevDebtNode.RemainingAmount
			}

			// * paid off
			if totalRemainingAmount > 0 &&
				currentAsset-totalRemainingAmount > 0 &&
				currentAsset-totalRemainingAmount > rec.TotalMonthlyPaymentDebt-debt.MonthlyPayment &&
				eligiblePaidOff[j] {

				fmt.Println("paid off here", j, getMonthAndYear(startDate, q))

				currentAsset = currentAsset - totalRemainingAmount
				currentRemainingDebt = 0
				isPaidOff = true

				if err := s.db.Debt.Update(s.db.GDB, map[string]interface{}{
					"forecast_paid_off_date": getMonthAndYear(startDate, q),
				}, debt.ID); err != nil {
					return nil
				}

				// * update next debt to be eligible paid off
				eligiblePaidOff[j+1] = true
			} else {
				currentAsset = currentAsset - debt.MonthlyPayment
				currentRemainingDebt = totalRemainingAmount - calculateDebtPaidEachMonth(debt.MonthlyPayment, debt.AnnualInterest)
			}

			if currentRemainingDebt > 0 {
				totalRemainingDebt = totalRemainingDebt + currentRemainingDebt
			}

			if prevDebtNode != nil && prevDebtNode.RemainingAmount <= 0 {
				currentAsset = currentAsset + debt.MonthlyPayment
			}

			// * append debt
			if currentRemainingDebt >= 0 {
				datasets = append(datasets, &model.LineChart{
					Group: debt.Name,
					Key:   getMonth(startDate, q),
					Debt:  roundFloat(currentRemainingDebt),
				})

				setDebtNodeToCache(cache, &model.DataDebtNode{
					SessionID:       rec.ID,
					NodeName:        fmt.Sprintf("%d", i),
					DebtID:          debt.ID,
					Index:           j,
					RemainingAmount: currentRemainingDebt,
					MonthlyPayment:  debt.MonthlyPayment,
					IsPaidOff:       isPaidOff,
				})
			}

		}

		// * calculate current node
		curNode = calculateNode(rec, q, fmt.Sprintf("%d", i),
			currentAsset,                   // * dynamic
			totalRemainingDebt,             // * dynamic
			isPaidAllDebt(eligiblePaidOff)) // * dynamic

		if prevNode != nil && !prevNode.IsAchivedEmergencyFund && curNode.IsAchivedEmergencyFund {
			events[0] = getMonthAndYear(startDate, q)
		}

		if prevNode != nil && !prevNode.IsAchivedRainydayFund && curNode.IsAchivedRainydayFund {
			events[1] = getMonthAndYear(startDate, q)
		}

		if prevNode != nil && !prevNode.IsAchivedInvestment && curNode.IsAchivedEmergencyFund && curNode.IsAchivedRainydayFund {
			events[2] = getMonthAndYear(startDate, q)
			monthlyNetFlowWithoutDebt = curNode.MonthlyNetFlow
			currentAssetToMillionaire = currentAsset
		}

		if prevNode != nil && !prevNode.IsAchivedFinancialFreedom && curNode.IsAchivedFinancialFreedom {
			events[3] = getMonthAndYear(startDate, q)
		}

		setNodeToCache(cache, curNode)

		// * append asset
		if curNode.CurrentAsset > 0 {
			datasets = append(datasets, &model.LineChart{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: curNode.CurrentAsset,
			})
		} else {
			datasets = append(datasets, &model.LineChart{
				Group: "Assets",
				Key:   getMonth(startDate, q),
				Asset: 0,
			})
			events[4] = getMonthAndYear(startDate, q)
			break
		}
	}

	if currentAssetToMillionaire > 0 {
		remainingAssetToMillionaire := MillionaireRate - currentAssetToMillionaire

		becomeMillionaireIn := remainingAssetToMillionaire / monthlyNetFlowWithoutDebt // months

		// fmt.Printf("becomeMillionaireIn ==== %f years \n", becomeMillionaireIn/12/2)

		t := startDate.AddDate(0, int(becomeMillionaireIn/2), 0)
		millionaireDate = t.Format("Jan 2006")
	} else {

	}

	if err := s.db.Session.Update(s.db.GDB, map[string]interface{}{
		"forecast_emergency_budget_filled_date": events[0],
		"forecast_start_investing_date":         events[2],
		"forecast_rainyday_budget_filled_date":  events[1],
		"forecast_financial_freedom_date":       events[3],
		"forecast_millionaire_date":             millionaireDate,
		"forecast_bankrupt":                     events[4],
	}, rec.ID); err != nil {
		fmt.Println("Error updating forecast events")
	}

	return datasets
}

func calculateDebtPaidEachMonth(monthlyPayment, annualInterest float64) float64 {
	return roundFloat((1 - (annualInterest / 12.0 / 100)) * monthlyPayment)
}

func calculateNode(session *model.Session, calcTime int64, nodeName string, currentAsset, totalRemainingDebt float64, isPaidAllDebt bool) *model.DataNode {
	var expectedEmergencyFund, expectedRainydayFund, expectedFunFund, actualEmergencyFund, actualRainydayFund, actualFunFund, retirementPlan float64
	var isAchivedInvestment, isAchivedEmergencyFund, isAchivedRainydayFund, isAchivedRetirementPlan bool

	isAchivedFinancialFreedom := true

	monthlyNetFlow := calculateMonthlyNetFlow(session)

	totalMonthlyPaymentDebt := session.TotalMonthlyPaymentDebt
	totalIncome := session.TotalAllIncome
	totalExpense := session.TotalAllExpense

	totalEssentialExpense := session.TotalEssentialExpense

	expectedEmergencyFund = roundFloat(totalEssentialExpense * EmergencyFundRate)
	expectedRainydayFund = roundFloat(totalEssentialExpense * RainydayFundRate)
	expectedFunFund = roundFloat(monthlyNetFlow * FunFundRate)

	retirementPlan = roundFloat(totalEssentialExpense * 12 * RetirementPlanRate)

	// * only achived when emergency fund and rainy day fund is achived and no debt
	if isPaidAllDebt {
		netAsset := currentAsset - totalRemainingDebt

		actualEmergencyFund := roundFloat(netAsset)
		if actualEmergencyFund <= 0 {
			actualEmergencyFund = 0
		}
		actualRainydayFund := roundFloat(netAsset - expectedEmergencyFund)
		if actualRainydayFund <= 0 {
			actualRainydayFund = 0
		}

		isAchivedEmergencyFund = netAsset >= expectedEmergencyFund
		if isAchivedEmergencyFund {
			actualEmergencyFund = expectedEmergencyFund
		}

		isAchivedRainydayFund = netAsset >= (expectedEmergencyFund + expectedRainydayFund)
		if isAchivedRainydayFund {
			actualRainydayFund = expectedRainydayFund
		}

		if isAchivedEmergencyFund && isAchivedRainydayFund {
			isAchivedInvestment = true

			// * calculate R
			r := currentAsset - (actualEmergencyFund + actualRainydayFund)

			if r >= 0 {
				currentAsset = currentAsset + (r * UserBudgetIncreasementRate)

				// * update forecast financial freedom
				if r*UserBudgetIncreasementRate < totalEssentialExpense {
					isAchivedFinancialFreedom = false
				}
			}

			isAchivedRetirementPlan = netAsset >= retirementPlan && retirementPlan > 0
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
		SessionID:                 session.ID,
		NodeName:                  nodeName,
		CurrentAsset:              roundFloat(currentAsset),
		TotalAllIncome:            totalIncome,
		TotalAllExpense:           totalExpense,
		TotalMonthlyPaymentDebt:   roundFloat(totalMonthlyPaymentDebt),
		MonthlyNetFlow:            monthlyNetFlow,
		Status:                    status,
		Descrtiption:              model.SessionStatusDescriptions[status],
		ExpectedEmergencyFund:     expectedEmergencyFund,
		ExpectedRainydayFund:      expectedRainydayFund,
		ExpectFunFund:             expectedFunFund,
		ActualEmergencyFund:       actualEmergencyFund,
		ActualRainydayFund:        actualRainydayFund,
		ActualFunFund:             actualFunFund,
		RetirementPlan:            retirementPlan,
		IsAchivedEmergencyFund:    isAchivedEmergencyFund,
		IsAchivedRainydayFund:     isAchivedRainydayFund,
		IsAchivedInvestment:       isAchivedInvestment,
		IsAchivedRetirementPlan:   isAchivedRetirementPlan,
		IsAchivedFinancialFreedom: isAchivedFinancialFreedom,
		IsPaidAllDebt:             isPaidAllDebt,
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

	cache.Set(fmt.Sprintf("%d_%s_expected_fun_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.ExpectFunFund)))
	cache.Set(fmt.Sprintf("%d_%s_actual_fun_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.ActualFunFund)))

	cache.Set(fmt.Sprintf("%d_%s_retirement_plan", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%f", node.RetirementPlan)))

	cache.Set(fmt.Sprintf("%d_%s_is_achived_emergency_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedEmergencyFund)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_rainyday_fund", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedRainydayFund)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_investment", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedInvestment)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_retirement_plan", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedRetirementPlan)))
	cache.Set(fmt.Sprintf("%d_%s_is_achived_financial_freedom", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsAchivedFinancialFreedom)))

	cache.Set(fmt.Sprintf("%d_%s_is_paid_all_debt", node.SessionID, node.NodeName), []byte(fmt.Sprintf("%v", node.IsPaidAllDebt)))
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

	ceff, _ := cache.Get(fmt.Sprintf("%d_%s_expected_fun_fund", session.ID, nodeName))
	cExpectFunFund, _ := strconv.ParseFloat(string(ceff), 64)

	caff, _ := cache.Get(fmt.Sprintf("%d_%s_actual_fun_fund", session.ID, nodeName))
	cActualFunFund, _ := strconv.ParseFloat(string(caff), 64)

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

	ciaff, _ := cache.Get(fmt.Sprintf("%d_%s_is_achived_financial_freedom", session.ID, nodeName))
	cIsAchivedFinancialFreedom, _ := strconv.ParseBool(string(ciaff))

	cipad, _ := cache.Get(fmt.Sprintf("%d_%s_is_paid_all_debt", session.ID, nodeName))
	cIsPaidAllDebt, _ := strconv.ParseBool(string(cipad))

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
		ExpectFunFund:         cExpectFunFund,
		ActualEmergencyFund:   cActualEmergencyFund,
		ActualRainydayFund:    cActualRainydayFund,
		ActualFunFund:         cActualFunFund,

		RetirementPlan: cRetirementPlan,

		IsAchivedEmergencyFund:    cIsAchivedEmergencyFund,
		IsAchivedRainydayFund:     cIsAchivedRainydayFund,
		IsAchivedInvestment:       cIsAchivedInvestment,
		IsAchivedRetirementPlan:   cIsAchivedRetirementPlan,
		IsAchivedFinancialFreedom: cIsAchivedFinancialFreedom,
		IsPaidAllDebt:             cIsPaidAllDebt,
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

func mappingFullStatus(status string) string {
	switch status {
	case model.SessionStatusBD:
		return "Budget Deficit"
	case model.SessionStatusPC2PC:
		return "Pay Check to Pay Check"
	case model.SessionStatusLFF:
		return "Limited financial flexibility"
	case model.SessionStatusGFF:
		return "Good financial flexibility"
	default:
		return "Default"
	}
}

func isPaidAllDebt(m map[int]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}

func calculateMonthlyNetFlow(session *model.Session) float64 {
	totalIncome := session.TotalAllIncome
	totalEssentialExpense := session.TotalEssentialExpense
	totalNonEssentialExpense := session.TotalNonEssentialExpense

	totalExpense := roundFloat(totalEssentialExpense + totalNonEssentialExpense)

	return totalIncome - totalExpense
}
