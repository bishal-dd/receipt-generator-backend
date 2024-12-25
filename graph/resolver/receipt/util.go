package receipt

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/bishal-dd/receipt-generator-backend/helper/paginationUtil"
)

func convertEdges(edges []*paginationUtil.Edge[*model.Receipt]) []*model.ReceiptEdge {
    dataEdges := make([]*model.ReceiptEdge, len(edges))
    for i, edge := range edges {
        dataEdges[i] = &model.ReceiptEdge{
            Cursor: edge.Cursor,
            Node:   edge.Node,
        }
    }
    return dataEdges
}


func searchUpdateInput(input model.UpdateReceipt) (map[string]interface{}) {
    updateInput := make(map[string]interface{})
    
    // Only add non-nil fields to the updateInput map
    if input.RecipientName != nil {
        updateInput["recipient_name"] = input.RecipientName
    }
    if input.RecipientEmail != nil {
        updateInput["recipient_email"] = input.RecipientEmail
    }
    if input.RecipientAddress != nil {
        updateInput["recipient_address"] = input.RecipientAddress
    }
    if input.RecipientPhone != nil {
        updateInput["recipient_phone"] = input.RecipientPhone
    }
    if input.ReceiptNo != nil {
        updateInput["receipt_no"] = input.ReceiptNo
    }
    if input.PaymentMethod != nil {
        updateInput["payment_method"] = input.PaymentMethod
    }
    if input.PaymentNote != nil {
        updateInput["payment_note"] = input.PaymentNote
    }
    if input.UserID != nil {
        updateInput["user_id"] = input.UserID
    }
    if input.TotalAmount != nil {
        updateInput["total_amount"] = input.TotalAmount
    }
    
    return updateInput
}

