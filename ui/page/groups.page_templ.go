// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"movie_night/ui/components"
)

func Groups(groups []components.Group) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container\"><div class=\"row\"><div class=\"col-10\"><div class=\"input-group\"><input class=\"form-control\" type=\"text\" placeholder=\"Search for group\"> <button class=\"btn btn-outline-secondary\">Search</button></div></div><div class=\"col-2\"><a href=\"/groups/create\"><button class=\"btn btn-success\">Create group</button></a></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(groups) > 0 {
			for _, group := range groups {
				templ_7745c5c3_Err = group.Render().Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div>You are not a part of any groups.</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
