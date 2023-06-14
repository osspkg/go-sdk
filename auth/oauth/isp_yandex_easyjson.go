// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package oauth

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD1fc6ea8DecodeGithubComosspkgGoSdkAuthOauth(in *jlexer.Lexer, out *modelYandex) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "display_name":
			out.Name = string(in.String())
		case "default_avatar_id":
			out.Icon = string(in.String())
		case "default_email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD1fc6ea8EncodeGithubComosspkgGoSdkAuthOauth(out *jwriter.Writer, in modelYandex) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"display_name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"default_avatar_id\":"
		out.RawString(prefix)
		out.String(string(in.Icon))
	}
	{
		const prefix string = ",\"default_email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v modelYandex) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD1fc6ea8EncodeGithubComosspkgGoSdkAuthOauth(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v modelYandex) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD1fc6ea8EncodeGithubComosspkgGoSdkAuthOauth(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *modelYandex) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD1fc6ea8DecodeGithubComosspkgGoSdkAuthOauth(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *modelYandex) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD1fc6ea8DecodeGithubComosspkgGoSdkAuthOauth(l, v)
}
