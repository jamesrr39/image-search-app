package imagesearch

type PersistedDescriptorWithMatchScore struct {
	MatchScore MatchScore
	Descriptor *PersistedImageDescriptor
}

func NewPersistedDescriptorWithMatchScore(matchScore MatchScore, descriptor *PersistedImageDescriptor) *PersistedDescriptorWithMatchScore {
	return &PersistedDescriptorWithMatchScore{matchScore, descriptor}
}
