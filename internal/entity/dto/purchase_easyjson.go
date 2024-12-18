// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

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

func easyjson2ad9a821DecodeGithubComGoParkMailRu20242BogoSortInternalEntityDto(in *jlexer.Lexer, out *PurchaseResponse) {
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
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "seller_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.SellerID).UnmarshalText(data))
			}
		case "customer_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CustomerID).UnmarshalText(data))
			}
		case "adverts":
			if in.IsNull() {
				in.Skip()
				out.Adverts = nil
			} else {
				in.Delim('[')
				if out.Adverts == nil {
					if !in.IsDelim(']') {
						out.Adverts = make([]PreviewAdvertCard, 0, 0)
					} else {
						out.Adverts = []PreviewAdvertCard{}
					}
				} else {
					out.Adverts = (out.Adverts)[:0]
				}
				for !in.IsDelim(']') {
					var v1 PreviewAdvertCard
					(v1).UnmarshalEasyJSON(in)
					out.Adverts = append(out.Adverts, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "address":
			out.Address = string(in.String())
		case "status":
			out.Status = PurchaseStatus(in.String())
		case "payment_method":
			out.PaymentMethod = PaymentMethod(in.String())
		case "delivery_method":
			out.DeliveryMethod = DeliveryMethod(in.String())
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
func easyjson2ad9a821EncodeGithubComGoParkMailRu20242BogoSortInternalEntityDto(out *jwriter.Writer, in PurchaseResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"seller_id\":"
		out.RawString(prefix)
		out.RawText((in.SellerID).MarshalText())
	}
	{
		const prefix string = ",\"customer_id\":"
		out.RawString(prefix)
		out.RawText((in.CustomerID).MarshalText())
	}
	{
		const prefix string = ",\"adverts\":"
		out.RawString(prefix)
		if in.Adverts == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Adverts {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix)
		out.String(string(in.Status))
	}
	{
		const prefix string = ",\"payment_method\":"
		out.RawString(prefix)
		out.String(string(in.PaymentMethod))
	}
	{
		const prefix string = ",\"delivery_method\":"
		out.RawString(prefix)
		out.String(string(in.DeliveryMethod))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PurchaseResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2ad9a821EncodeGithubComGoParkMailRu20242BogoSortInternalEntityDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PurchaseResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2ad9a821EncodeGithubComGoParkMailRu20242BogoSortInternalEntityDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PurchaseResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2ad9a821DecodeGithubComGoParkMailRu20242BogoSortInternalEntityDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PurchaseResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2ad9a821DecodeGithubComGoParkMailRu20242BogoSortInternalEntityDto(l, v)
}
func easyjson2ad9a821DecodeGithubComGoParkMailRu20242BogoSortInternalEntityDto1(in *jlexer.Lexer, out *PurchaseRequest) {
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
		case "cart_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CartID).UnmarshalText(data))
			}
		case "address":
			out.Address = string(in.String())
		case "payment_method":
			out.PaymentMethod = PaymentMethod(in.String())
		case "delivery_method":
			out.DeliveryMethod = DeliveryMethod(in.String())
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
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
func easyjson2ad9a821EncodeGithubComGoParkMailRu20242BogoSortInternalEntityDto1(out *jwriter.Writer, in PurchaseRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cart_id\":"
		out.RawString(prefix[1:])
		out.RawText((in.CartID).MarshalText())
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"payment_method\":"
		out.RawString(prefix)
		out.String(string(in.PaymentMethod))
	}
	{
		const prefix string = ",\"delivery_method\":"
		out.RawString(prefix)
		out.String(string(in.DeliveryMethod))
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserID).MarshalText())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PurchaseRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2ad9a821EncodeGithubComGoParkMailRu20242BogoSortInternalEntityDto1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PurchaseRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2ad9a821EncodeGithubComGoParkMailRu20242BogoSortInternalEntityDto1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PurchaseRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2ad9a821DecodeGithubComGoParkMailRu20242BogoSortInternalEntityDto1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PurchaseRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2ad9a821DecodeGithubComGoParkMailRu20242BogoSortInternalEntityDto1(l, v)
}
