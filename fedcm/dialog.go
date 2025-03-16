package fedcm

// DialogType represents the type of FedCM dialog
type DialogType string

const (
	// DialogTypeAccountList represents an account chooser dialog
	DialogTypeAccountList DialogType = "AccountChooser"
	// DialogTypeAutoReauth represents an auto reauthorization dialog
	DialogTypeAutoReauth DialogType = "AutoReauthn"
)

// Driver represents the interface that a WebDriver must implement to support FedCM operations
type Driver interface {
	GetDialogType() string
	GetTitle() string
	GetSubtitle() map[string]string
	GetAccountList() []map[string]string
	SelectAccount(index int) error
	Accept() error
	Dismiss() error
}

// Dialog represents a FedCM dialog that can be interacted with
type Dialog struct {
	driver Driver
}

// NewDialog creates a new Dialog instance
func NewDialog(driver Driver) *Dialog {
	return &Dialog{
		driver: driver,
	}
}

// GetType returns the type of the dialog currently being shown
func (d *Dialog) GetType() DialogType {
	return DialogType(d.driver.GetDialogType())
}

// GetTitle returns the title of the dialog
func (d *Dialog) GetTitle() string {
	return d.driver.GetTitle()
}

// GetSubtitle returns the subtitle of the dialog
func (d *Dialog) GetSubtitle() string {
	result := d.driver.GetSubtitle()
	if subtitle, ok := result["subtitle"]; ok {
		return subtitle
	}
	return ""
}

// GetAccounts returns the list of accounts shown in the dialog
func (d *Dialog) GetAccounts() []*Account {
	accountsData := d.driver.GetAccountList()
	accounts := make([]*Account, len(accountsData))
	for i, data := range accountsData {
		accounts[i] = NewAccount(data)
	}
	return accounts
}

// SelectAccount selects an account from the dialog by index
func (d *Dialog) SelectAccount(index int) error {
	return d.driver.SelectAccount(index)
}

// Accept clicks the continue button in the dialog
func (d *Dialog) Accept() error {
	return d.driver.Accept()
}

// Dismiss cancels/dismisses the dialog
func (d *Dialog) Dismiss() error {
	return d.driver.Dismiss()
}
