package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/intake"

type NewItemInput = intake.NewItemInput
type Draft = intake.Draft

func Slugify(title string) string      { return intake.Slugify(title) }
func ValidateTitle(title string) error { return intake.ValidateTitle(title) }
func InsertActiveRow(backlogContent, row string) (string, error) {
	return intake.InsertActiveRow(backlogContent, row)
}
func InsertNewItem(backlog string, item NewItemInput, date string) (string, error) {
	return intake.InsertNewItem(backlog, item, date)
}
