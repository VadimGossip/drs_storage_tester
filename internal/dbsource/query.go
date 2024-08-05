package dbsource

const sqlReadyRequestsQuery string = `
	select c.src_gwgr_id
		  ,c.anumber_in
		  ,c.bnumber_in
		  ,c.anumber_out
		  ,c.bnumber_out
		  ,c.ip
		  ,c.prefix
	  from cdr_based_drs_test_data c
`

const sqlReadyRequestsLimitQuery string = `
	select c.src_gwgr_id
		  ,c.anumber_in
		  ,c.bnumber_in
		  ,c.anumber_out
		  ,c.bnumber_out
		  ,c.ip
		  ,c.prefix
	  from cdr_based_drs_test_data c
	 where rownum <= %d 
`
