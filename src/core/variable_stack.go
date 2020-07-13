package core


type VarScope struct {
	super *Variables
	local *Variables
}

func newVarScope(super, local *Variables) *VarScope {
	return &VarScope{
		super: super,
		local: local,
	}
}
