use templategoREST;
db.createUser({ user: 'templategoRESTUser', pwd: 'Sf17a033vcF!',
		roles: [ { role: "readWrite", db: "templategoREST" } ]
	});
exit;
