INSERT INTO services (id, name, secret) 
VALUES ('512db16d-6d5b-4af4-aedd-a86e5425df30', 'test', 'test-secret')
ON CONFLICT DO NOTHING;