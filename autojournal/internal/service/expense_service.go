package service

import "autojournal/internal/domain"

type ExpenseService struct {
	repo domain.ExpenseRepository
}

type ExpenseCategorySevice struct {
	repo domain.ExpenseCategoryRepository
}

func NewExpense(repo domain.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func NewExpenseCategory(repo domain.ExpenseCategoryRepository) *ExpenseCategorySevice {
	return &ExpenseCategorySevice{repo: repo}
}
