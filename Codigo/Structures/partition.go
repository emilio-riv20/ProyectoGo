package structures

type PARTITION struct {
	PartStatus      [1]byte
	PartType        [1]byte
	PartFit         [1]byte
	PartStart       int32
	PartSize        int32
	PartName        [16]byte
	PartCorrelative int32
	PartId          [4]byte
}

func (p *PARTITION) CrearP(partStart, partSize int, partType, partFit, partName string) {
	p.PartStatus[0] = '0'

	p.PartStart = int32(partStart)

	p.PartSize = int32(partSize)

	if len(partType) > 0 {
		p.PartType[0] = partType[0]
	}

	if len(partFit) > 0 {
		p.PartFit[0] = partFit[0]
	}

	copy(p.PartName[:], partName)
}

func (p *PARTITION) MontarP(correlative int, id string) error {
	p.PartCorrelative = int32(correlative)
	copy(p.PartId[:], id)

	return nil
}
