digraph finite_state_automaton {
  rankdir = LR;
  fontname="sans-serif";
  node [shape=circle, fontname="sans-serif"];
  edge [fontname="sans-serif"];
{{if .Transition}}{{printf  "  // transition\n"}}{{end -}}
{{range $k, $v := .Transition -}}
  {{"  "}}"{{$k.State}}" -> "{{$v}}" [label="{{$k.Input}}"];
{{end -}}
{{"  "}}// start
{{"  "}}start [shape=point];
{{"  "}}start -> "{{.Start}}";
{{if .FinalState}}{{printf "  // final states\n"}}{{end -}}
{{range $k, $v := .FinalState -}}
  {{"  "}}"{{$k}}" [shape=ellipse, peripheries=2];
{{end -}}
}
