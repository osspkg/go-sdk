// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jwt

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

func easyjson171edd05DecodeGithubComosspkgGoSdkAuthJwt(in *jlexer.Lexer, out *Header) {
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
		case "kid":
			out.Kid = string(in.String())
		case "alg":
			out.Alg = string(in.String())
		case "iat":
			out.IssuedAt = int64(in.Int64())
		case "eat":
			out.ExpiresAt = int64(in.Int64())
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
func easyjson171edd05EncodeGithubComosspkgGoSdkAuthJwt(out *jwriter.Writer, in Header) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"kid\":"
		out.RawString(prefix[1:])
		out.String(string(in.Kid))
	}
	{
		const prefix string = ",\"alg\":"
		out.RawString(prefix)
		out.String(string(in.Alg))
	}
	{
		const prefix string = ",\"iat\":"
		out.RawString(prefix)
		out.Int64(int64(in.IssuedAt))
	}
	{
		const prefix string = ",\"eat\":"
		out.RawString(prefix)
		out.Int64(int64(in.ExpiresAt))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Header) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson171edd05EncodeGithubComosspkgGoSdkAuthJwt(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Header) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson171edd05EncodeGithubComosspkgGoSdkAuthJwt(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Header) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson171edd05DecodeGithubComosspkgGoSdkAuthJwt(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Header) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson171edd05DecodeGithubComosspkgGoSdkAuthJwt(l, v)
}
