package lib

import (
	db "geo-git/lib/database"
	"path/filepath"
)

type Respository struct {
	GitPath   string
	Database  *db.Database
	Index     *Index
	Refs      *Refs
	Workspace *Workspace
}

func NewRepository(git_path string) *Respository {
	db := db.NewDatabase(filepath.Join(git_path, "objects"))
	index := NewIndex(filepath.Join(git_path, "index"))
	refs := NewRefs(git_path)
	workspace := NewWorkSpace(filepath.Dir(git_path))
	return &Respository{GitPath: git_path, Database: db, Index: index,
		Refs: refs, Workspace: workspace}
}

func (r *Respository) GetDatabase() *db.Database {
	return r.Database
}

func (r *Respository) GetIndex() *Index {
	return r.Index
}
func (r *Respository) GetRefs() *Index {
	return r.Index
}
func (r *Respository) GetWorkspace() *Workspace {
	return r.Workspace
}
