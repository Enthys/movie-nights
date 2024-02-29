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

func ViewGroup() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container\"><div class=\"row\"><img src=\"/assets/images/filler.jpeg\" class=\"col-sm-3 col-md-3 col-12\" style=\"max-height: 200px; object-fit: cover\"><div class=\"col-sm-9 col-md-9 col-12 \"><h3 class=\"\">This group has a length!</h3><p>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type  and scrambled it to make a type specimen book. It has survive</p></div></div><hr><div class=\"row mb-2\"><div class=\"col-sm-12 col-lg-4\"><div class=\"input-group\"><input type=\"text\" class=\"form-control\" placeholder=\"Movie IMDb link\"> <button type=\"button\" class=\"btn btn-outline-success\">Add</button></div></div><div class=\"col-sm-12 col-lg-1\"><p class=\"my-1 lh-5 text-center\"><span class=\"align-baseline\">Or</span></p></div><div class=\"col-sm-12 col-lg-4\"><div class=\"input-group\"><input type=\"text\" class=\"form-control\" placeholder=\"Movie name\"> <button type=\"button\" class=\"btn btn-outline-success\">Seach</button></div></div><div class=\"col-sm-12 col-lg-1\"><p class=\"my-1 lh-5 text-center\"><span class=\"align-baseline\">Or</span></p></div><div class=\"px-2 px-sm-0 col-12 col-sm-2 text-center\"><button class=\" btn btn-success\">Pick Random</button></div></div><div class=\"row\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i := 0; i < 10; i++ {
			templ_7745c5c3_Err = components.NewMovie("foo", "foo", "foo", []string{}).Render().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}