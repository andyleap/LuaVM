t = {}
t.a = 1
t.b = 2

function t:test()
	print(t.a + t.b)
end

function test(...)
	print(...)
	t = {...}
	print(t[2])
end

test(t.a, t.b)
