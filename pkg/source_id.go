package provider

type sourceIDType int32

const (
	sourceIDOpenFeatureGo sourceIDType = 103
)

func (s sourceIDType) Int32() int32 {
	return int32(s)
}
