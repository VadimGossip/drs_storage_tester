package dbsource

const sqlSUPGWQuery = `
     select distinct g.gw_gwgr_id    
      from gateways g
      join prefixes p on g.gw_id = p.pfx_gw_id
     where g.gw_orig_ability = 1
       and p.pfx_direction = 'originate'`
