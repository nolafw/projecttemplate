RDBで常に親子関係が成立するリレーションのentityの場合

```
@startuml Orders_OrderDetails_Relationship

entity "orders" as orders {
  * id : UUID <<PK>>
  --
  * user_id : UUID <<FK>>
  * created_at : timestamp
  * updated_at : timestamp
  status : string
  total_amount : decimal
  payment_method : string
  shipping_address : string
}

entity "order_details" as orderDetails {
  * id : UUID <<PK>>
  --
  * order_id : UUID <<FK>>
  * product_id : UUID <<FK>>
  * quantity : integer
  * price : decimal
  * created_at : timestamp
  * updated_at : timestamp
  options : json
}

orders ||--o{ orderDetails : "has many"

@enduml
```