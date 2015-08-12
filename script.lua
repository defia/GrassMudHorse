function score()
		--print(droprate(),averagelatency(),address(),((1-droprate())^50)/ averagelatency())

	    return ((1-droprate())^50)/ averagelatency()--this is actually the native score function translated in lua.
end
