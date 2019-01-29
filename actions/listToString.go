package actions

import (
	"encoding/json"
	"fmt"
	"tgsd96/onend/app"
)

type QueryResult struct {
	CustID int64  `json:"cust_id"`
	Col    string `json:"col"`
	Marg   string `json:"marg"`
	Total  string `json:"total"`
	Name   string `json:"name"`
	Area   string `json:"area"`
}

func ExecuteListSQL(startDate string, endDate string) []byte {
	sqlQuery := `
	SELECT A.*, B.Name, B.area FROM (
	SELECT COALESCE(A.CUST_ID, B.CUST_ID) as CUST_ID,
	A.COL, A.MARG, B.Total from
	(SELECT * 
	FROM crosstab(
	'select  cust_id,interface_code,string_agg((amount::numeric)::text,''|'')
	FROM PURCHASES 
	where created_at >= ''%s''
	and created_at <= ''%s''
	group by cust_id,interface_code order by 1,2','
	SELECT DISTINCT INTERFACE_CODE from masters order by 1')
	AS ct(cust_id int, COL text, MARG text))A
	LEFT OUTER JOIN
		(select cust_id, sum(amount) as Total 
			from ledgers group by cust_id 
				having sum(amount)>0) B 
	ON A.cust_id = B.cust_id) A 
	INNER JOIN 
	master_view B ON a.cust_id = B.cust_id`

	stringQuery := fmt.Sprintf(sqlQuery, startDate, endDate)

	fmt.Println(stringQuery)
	// declare the output structure
	var queryResult []QueryResult

	app.App.DB.LogMode(true)

	// execute the query
	if err := app.App.DB.Raw(stringQuery).Scan(&queryResult).Error; err != nil {
		fmt.Println(err)
		return []byte("")
	}

	msg, _ := json.Marshal(&queryResult)
	return msg
}
