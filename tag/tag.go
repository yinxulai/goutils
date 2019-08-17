package tag

type TagTree struct {
	Key  string
	Type string
	Tags map[string][]string
}

func Read(target interface{}) *TagTree {
	return nil
}
