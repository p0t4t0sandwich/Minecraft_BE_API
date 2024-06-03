// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Root() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Sci-Fi Minecraft Camp Administration</title></head><script>\n\n\t\t// Label handling\n\t\tlet label = \"default\";\n\t\tfunction setLabel() {\n\t\t\tlabel = document.getElementById(\"serverLabel\").value;\n\t\t}\n\t\twindow.onload = setLabel;\n\n\t\t// Command sending\n\t\tfunction sendCmd(cmd) {\n\t\t\tfetch(`/api/cmd/${label}`, {\n\t\t\t\tmethod: \"POST\",\n\t\t\t\tbody: JSON.stringify({cmd: cmd}),\n\t\t\t\theaders: {\"Content-Type\": \"application/json\"}\n\t\t\t});\n\t\t}\n\n        function fiveMinuteWarning() {\n            sendCmd(\"title @a times 0 100 0\");\n\t\t\tsendCmd(\"title @a title 5 minute warning!\");\n        }\n\n\t\tfunction flashDarkness() {\n\t\t\tsendCmd(\"effect @a blindness 2 2\");\n\t\t}\n\n\t\t// Entity clearing forever -- TODO: hardcoded for now\n\t\tentityList = [\n\t\t\t\"ender_dragon\",\n\t\t\t\"warden\",\n\t\t\t\"wither\",\n\t\t\t\"lingering_potion\"\n\t\t]\n\t\t\n\t\tsetInterval(() => {\n\t\t\tentityList.forEach(entity => sendCmd(`kill @e[type=${entity}]`));\n\t\t}, 1000);\n    </script><body><!-- This is a dummy frame to prevent the page from reloading when a form is submitted --><iframe name=\"dummyframe\" id=\"dummyframe\" style=\"display: none;\"></iframe><h1>Sci-Fi Minecraft Camp Administration</h1><!-- Set server label --><p>Set server label:</p><form onSubmit=\"setLabel()\" target=\"dummyframe\"><input type=\"text\" id=\"serverLabel\" name=\"serverLabel\" placeholder=\"Server Label\"> <input type=\"submit\" value=\"Submit\"></form><!-- Pause buttons --><p>Pause buttons:</p><button onclick=\"sendCmd(&#39;globalpause true&#39;)\">Pause Server</button> <button onclick=\"sendCmd(&#39;globalpause false&#39;)\">Unpause Server</button><!-- Attention Utils --><p>Attention Utils:</p><button onclick=\"fiveMinuteWarning()\">Five Min Warning</button> <button onclick=\"flashDarkness()\">Flash Darkness</button><br><p>General Utils:</p><p>Clear player inventory:</p><form onSubmit=\"sendCmd(&#39;clear &#39; + document.getElementById(&#39;clearPlayer&#39;).value + &#39; &#39; + document.getElementById(&#39;clearItem&#39;).value)\" target=\"dummyframe\"><input type=\"text\" id=\"clearPlayer\" name=\"player\" placeholder=\"Player\"> <input type=\"text\" id=\"clearItem\" name=\"item\" placeholder=\"Item\"> <input type=\"submit\" value=\"Clear\"></form><p>Remove entity:</p><form onSubmit=\"sendCmd(&#39;kill @e[type=&#39; + document.getElementById(&#39;entityType&#39;).value + &#39;]&#39;)\" target=\"dummyframe\"><input type=\"text\" id=\"entityType\" name=\"entityType\" placeholder=\"Entity Type\"> <input type=\"submit\" value=\"Remove\"></form><!-- Custom command --><p>Custom command:</p><form onSubmit=\"sendCmd(document.getElementById(&#39;customCmd&#39;).value)\" target=\"dummyframe\"><input type=\"text\" id=\"customCmd\" name=\"customCmd\" placeholder=\"Custom Command\"> <input type=\"submit\" value=\"Submit\"></form></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}