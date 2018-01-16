--cfg--
box.cfg{listen = 1311}

--Space--
box.schema.create_space('tnt_space')
box.space.tnt_space:create_index('primary')
box.space.tnt_space:insert{1, "foo"}
box.space.tnt_space:insert{2, "bar"}

--Space with ns--
box.schema.create_space('ns_tnt_space')
box.space.ns_tnt_space:create_index('primary')
box.space.ns_tnt_space:insert{1, "qux"}

--User--
box.schema.user.create('tnt_test', {password = 'tnt_test'})
box.schema.user.grant('tnt_test', 'read,write,execute', 'universe')
