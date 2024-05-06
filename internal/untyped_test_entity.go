package internal

import (
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UntypedTestEntity struct {
	additionalData map[string]interface{}
	id             *string
	title          *string
	location       absser.UntypedNodeable
	keywords       absser.UntypedNodeable
	detail         absser.UntypedNodeable
	table          absser.UntypedNodeable
}

type TestUntypedTestEntityable interface {
	absser.Parsable
	absser.AdditionalDataHolder
	GetId() *string
	SetId(value *string)
	GetTitle() *string
	SetTitle(value *string)
	GetLocation() absser.UntypedNodeable
	SetLocation(value absser.UntypedNodeable)
	GetKeywords() absser.UntypedNodeable
	SetKeywords(value absser.UntypedNodeable)
	GetDetail() absser.UntypedNodeable
	SetDetail(value absser.UntypedNodeable)
	GetTable() absser.UntypedNodeable
	SetTable(value absser.UntypedNodeable)
}

func NewUntypedTestEntity() *UntypedTestEntity {
	return &UntypedTestEntity{
		additionalData: make(map[string]interface{}),
	}
}

func UntypedTestEntityDiscriminator(parseNode absser.ParseNode) (absser.Parsable, error) {
	return NewUntypedTestEntity(), nil
}

func (e *UntypedTestEntity) GetAdditionalData() map[string]interface{} {
	return e.additionalData
}
func (e *UntypedTestEntity) SetAdditionalData(value map[string]interface{}) {
	e.additionalData = value
}

func (e *UntypedTestEntity) GetId() *string {
	return e.id
}
func (e *UntypedTestEntity) SetId(value *string) {
	e.id = value
}

func (e *UntypedTestEntity) GetTitle() *string {
	return e.title
}
func (e *UntypedTestEntity) SetTitle(value *string) {
	e.title = value
}

func (e *UntypedTestEntity) GetLocation() absser.UntypedNodeable {
	return e.location
}
func (e *UntypedTestEntity) SetLocation(value absser.UntypedNodeable) {
	e.location = value
}

func (e *UntypedTestEntity) GetKeywords() absser.UntypedNodeable {
	return e.keywords
}
func (e *UntypedTestEntity) SetKeywords(value absser.UntypedNodeable) {
	e.keywords = value
}

func (e *UntypedTestEntity) GetDetail() absser.UntypedNodeable {
	return e.detail
}
func (e *UntypedTestEntity) SetDetail(value absser.UntypedNodeable) {
	e.detail = value
}

func (e *UntypedTestEntity) GetTable() absser.UntypedNodeable {
	return e.table
}
func (e *UntypedTestEntity) SetTable(value absser.UntypedNodeable) {
	e.table = value
}

func (e *UntypedTestEntity) GetFieldDeserializers() map[string]func(absser.ParseNode) error {
	res := make(map[string]func(absser.ParseNode) error)
	res["id"] = func(n absser.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetId(val)
		}
		return nil
	}
	res["title"] = func(n absser.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			e.SetTitle(val)
		}
		return nil
	}
	res["location"] = func(n absser.ParseNode) error {
		val, err := n.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		if err != nil {
			return err
		}
		if val != nil {
			e.SetLocation(val.(absser.UntypedNodeable))
		}
		return nil
	}
	res["keywords"] = func(n absser.ParseNode) error {
		val, err := n.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		if err != nil {
			return err
		}
		if val != nil {
			e.SetKeywords(val.(absser.UntypedNodeable))
		}
		return nil
	}
	res["detail"] = func(n absser.ParseNode) error {
		val, err := n.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		if err != nil {
			return err
		}
		if val != nil {
			e.SetDetail(val.(absser.UntypedNodeable))
		}
		return nil
	}
	res["table"] = func(n absser.ParseNode) error {
		val, err := n.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
		if err != nil {
			return err
		}
		if val != nil {
			e.SetTable(val.(absser.UntypedNodeable))
		}
		return nil
	}
	return res
}

func (m *UntypedTestEntity) Serialize(writer absser.SerializationWriter) error {
	{
		err := writer.WriteStringValue("id", m.GetId())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("title", m.GetTitle())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("location", m.GetLocation())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("keywords", m.GetKeywords())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("detail", m.GetDetail())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("table", m.GetTable())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteAdditionalData(m.GetAdditionalData())
		if err != nil {
			return err
		}
	}
	return nil
}
