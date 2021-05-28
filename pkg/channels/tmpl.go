package channels

import (
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
)

// parseTemplates parses added, deleted and updated templates
func parseTemplates(addedTmpl, deletedTmpl, updatedTmpl string) (
	tmplAdded, tmplDeleted, tmplUpdated *template.Template, err error) {
	funcMap := template.FuncMap{
		"field": func(s *unstructured.Unstructured, path ...string) string {
			// methods can be used in template:
			// s.GetName()
			// s.GetNamespace()
			str, exist, err := unstructured.NestedString(s.Object, path...)
			if err != nil {
				return fmt.Sprintf("[Error reading field .%s: %s]", strings.Join(path, "."), err)
			}
			if !exist {
				return fmt.Sprintf("[Field .%s not exist]", strings.Join(path, "."))
			}
			return str
		},
	}

	// example of using field:
	//tmpl = "[{{.Obj.GroupVersionKind}}] is created: " +
	// "{{.Obj.GetNamespace}}/{{.Obj.GetName}} {{field .Obj \"kind\"}}\n"
	if addedTmpl == "" {
		addedTmpl = "Resource [{{.Obj.GroupVersionKind}}, {{.Obj.GetNamespace}}/{{.Obj.GetName}}] has been added\n"
	}
	if deletedTmpl == "" {
		deletedTmpl = "Resource [{{.Obj.GroupVersionKind}}, {{.Obj.GetNamespace}}/{{.Obj.GetName}}] has been deleted\n"
	}
	if updatedTmpl == "" {
		updatedTmpl = "Resource [{{.Obj.GroupVersionKind}}, {{.Obj.GetNamespace}}/{{.Obj.GetName}}] has been updated\n"
	}

	tmplAdded, err = template.New("").Funcs(funcMap).Parse(addedTmpl)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "invalid added template")
	}

	tmplDeleted, err = template.New("").Funcs(funcMap).Parse(deletedTmpl)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "invalid  deleted template")
	}

	tmplUpdated, err = template.New("").Funcs(funcMap).Parse(updatedTmpl)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "invalid  updated template")
	}

	return tmplAdded, tmplDeleted, tmplUpdated, nil
}
