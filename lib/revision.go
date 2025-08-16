package lib

import "regexp"

const PARENT = `/^(.+)\^$/`
const ANCESTOR = `/^(.+)~(\d+)$/`

var REF_ALIASES = map[string]string{
	"@": "HEAD",
}

type Revision struct {
	ref      *RevisionRef
	paren    *RevisionParent
	ancestor *RevisionAncestor
}

type RevisionRef struct {
}

type RevisionParent struct {
}

type RevisionAncestor struct {
}

func NewRevision() *Revision {
	return &Revision{}
}

func Parse(revision string) (*Revision, error) {
	//TODO: WRite this function in a way to use multiple patterns
	// re := regexp.MustCompile(pattern)
	// re.MatchString(string)
	//then submatches := re.FindStringSubmatch(input)
	matched, err := regexp.MatchString(PARENT, revision)
	if err != nil {
		return nil, err
	}

	if matched {
		// match, _ := regexp.Match(PARENT, []byte(revision))
		// rev := Parse(mat)
	}
	return nil, nil
}
