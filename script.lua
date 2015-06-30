function score()
		print(droprate(),averagelatency(),address(),((1-droprate())^50)/ averagelatency())
	    return ((1-droprate())^50)/ averagelatency()
end
