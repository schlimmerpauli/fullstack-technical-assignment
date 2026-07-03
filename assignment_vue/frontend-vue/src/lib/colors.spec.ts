import { describe, expect, it } from 'vitest'

import { classicColorKeys, filterClassicColors } from './colors'

describe('colors helpers', () => {
	it('keeps only the classic five colors in their original order', () => {
		expect(
			filterClassicColors(['silver', 'blue', 'red', 'pink', 'green', 'white', 'black', 'orange', 'blue']),
		).toEqual(['blue', 'red', 'green', 'white', 'black'])
	})

	it('defines the classic five color keys', () => {
		expect(classicColorKeys).toEqual(['black', 'blue', 'green', 'red', 'white'])
	})
})