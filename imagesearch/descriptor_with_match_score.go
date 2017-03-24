package imagesearch

type DescriptorWithMatchScore struct {
	MatchScore MatchScore
	Descriptor *ImageDescriptor
}

func NewDescriptorWithMatchScore(matchScore MatchScore, descriptor *ImageDescriptor) *DescriptorWithMatchScore {
	return &DescriptorWithMatchScore{matchScore, descriptor}
}
