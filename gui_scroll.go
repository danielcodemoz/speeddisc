package main

// scrollThumb is independent of the OS so the proportional geometry and drag
// rounding can be tested without a desktop. A range maximum is inclusive.
func scrollThumb(height, minThumb, low, high, page, pos int32) (start, size, last int32) {
	rng := int64(high) - int64(low) + 1
	if page < 1 {
		page = 1
	}
	if height < 1 || rng <= int64(page) {
		return 0, height, low
	}
	extent := rng - int64(page)
	last = int32(int64(low) + extent)
	size = min(height, max(minThumb, int32(int64(height)*int64(page)/rng)))
	travel := int64(height - size)
	step := int64(min(last, max(low, pos))) - int64(low)
	return int32((step*travel + extent/2) / extent), size, last
}

func scrollTarget(start, travel, low, last int32) int32 {
	if travel <= 0 || last <= low {
		return low
	}
	start = min(travel, max(0, start))
	return int32(int64(low) + (int64(start)*(int64(last)-int64(low))+int64(travel)/2)/int64(travel))
}
