package newcontext

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
)

func newForm() (*huh.Form, huh.Field, huh.Field, *huh.Confirm) {
	nameInput := huh.NewInput().
		Key("name").
		Title("Human-readable name:").
		Placeholder("AWS Prod RDS")
	nameInput.Focus()

	dsnInput := huh.NewInput().
		Key("dsn").
		Title("DSN").
		Placeholder("mysql://root:notrealpassword@(0.0.0.0:3306)/test")

	confirm := huh.NewConfirm().
		Key("done").
		Title("Are you sure?").
		Affirmative("Yes").
		Negative("No")

	form := huh.NewForm(
		huh.NewGroup(
			nameInput,
			dsnInput,
		),
		huh.NewGroup(
			confirm,
		),
	).WithTheme(newHuhTheme()).WithShowHelp(false)

	return form, nameInput, dsnInput, confirm
}

func newHuhTheme() *huh.Theme {
	return &huh.Theme{
		Form: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(color.Border),
		Group: lipgloss.NewStyle().
			Align(lipgloss.Center),
		Blurred: huh.FieldStyles{
			TextInput: newTextInputStyles(),
		},
		Focused: huh.FieldStyles{
			BlurredButton: lipgloss.NewStyle().
				Foreground(color.Text).
				Padding(0, 1).
				Margin(1),
			FocusedButton: lipgloss.NewStyle().
				Foreground(color.Text).
				Background(color.MainAccent).
				Padding(0, 1).
				Margin(1),
			TextInput: newTextInputStyles(),
		},
	}
}

func newTextInputStyles() huh.TextInputStyles {
	return huh.TextInputStyles{
		Placeholder: lipgloss.NewStyle().Foreground(color.Placeholder),
		Text:        lipgloss.NewStyle().Foreground(color.Text),
		Cursor:      lipgloss.NewStyle().Foreground(color.Text),
		Prompt:      lipgloss.NewStyle().Foreground(color.MainAccent),
	}
}
