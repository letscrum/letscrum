package note

import (
	"fmt"
	"os"
	"strings"

	"github.com/letscrum/letscrum/tools/releasenoter/utils"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

var (
	noteKindLevel = map[string]int{"bug-fix": 1, "security-fix": 2, "feature": 3, "test": 4, "other": 5}
	noteAreaLevel = map[string]int{"apis": 1, "architecture": 2, "infrastructure": 3, "installation": 4, "documentation": 5, "others": 6}
)

type Note struct {
	Name  string   `json:"name"`
	Kind  string   `json:"kind"`
	Area  string   `json:"area"`
	Notes []string `json:"notes"`
}

func ParseNotesFile(files []string) ([]Note, error) {
	var notes []Note
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return nil, errors.Errorf("could not find file, %s does not exist", file)
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return nil, errors.Errorf("unable to open file %s: %s", file, err.Error())
		}

		var note Note
		if err = yaml.Unmarshal(content, &note); err != nil {
			return nil, errors.Errorf("unable to parse release note %s:%s", file, err.Error())
		}
		note.Name = file
		notes = append(notes, note)
	}

	if len(notes) < 1 {
		return nil, errors.New("failed to find any release notes")
	}

	return notes, nil
}

func CreateMarkDown(notes []Note, outPath, version string) error {
	kindMap, err := getKindMap(notes)
	if err != nil {
		return errors.Wrap(err, "get kind map err")
	}

	content := fmt.Sprintf("# release  %s\n", version)
	for _, kind := range kindMap.Array {
		kindNotes := kindMap.Map[kind]
		areaNotes, errK := getAreaMap(kindNotes.([]Note))
		if errK != nil {
			return errors.Wrap(errK, "get area map err")
		}
		content += joinContent(kind, areaNotes)
	}

	if !strings.HasSuffix(outPath, "/") {
		outPath += "/"
	}
	err = os.WriteFile(outPath+version+".md", []byte(content), 0o666)
	if err != nil {
		return errors.Errorf("create markdown file err:%s", err.Error())
	}

	return nil
}

func getKindMap(notes []Note) (utils.OrderMap, error) {
	kindMap := utils.OrderMap{Map: map[string]interface{}{}}
	for _, note := range notes {
		if _, ok := noteKindLevel[note.Kind]; !ok {
			return utils.OrderMap{}, errors.Errorf("%s not present kind name is %s", note.Name, note.Kind)
		}
		kindMap.Sort(note.Kind, noteKindLevel)
		if kindMap.Map[note.Kind] == nil {
			kindMap.Map[note.Kind] = []Note{}
		}
		kindMap.Map[note.Kind] = append(kindMap.Map[note.Kind].([]Note), note)
	}

	return kindMap, nil
}

func getAreaMap(notes []Note) (utils.OrderMap, error) {
	areaMap := utils.OrderMap{Map: map[string]interface{}{}}
	for _, note := range notes {
		if _, ok := noteAreaLevel[note.Area]; !ok {
			return utils.OrderMap{}, errors.Errorf("%s not present area name is %s", note.Name, note.Area)
		}
		areaMap.Sort(note.Area, noteAreaLevel)
		if areaMap.Map[note.Area] == nil {
			areaMap.Map[note.Area] = []string{}
		}
		areaMap.Map[note.Area] = append(areaMap.Map[note.Area].([]string), editNoteFormat(note.Notes)...)
	}
	return areaMap, nil
}

func editNoteFormat(notes []string) []string {
	for i := range notes {
		notes[i] = "\n+ " + notes[i]
	}
	return notes
}

func joinContent(kind string, notes utils.OrderMap) string {
	content := fmt.Sprintf("## %s\n", kind)
	for _, note := range notes.Array {
		content += fmt.Sprintf("#### %s", note)
		content += strings.Join(notes.Map[note].([]string), "")
		content += "\n"
	}
	return content
}
