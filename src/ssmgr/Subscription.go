package ssmgr

// SubscriptionList is the list of subscribed sources.
type SubscriptionList []string

// This returns true self.
func (sl *SubscriptionList) This() []string {
	return []string(*sl)
}
