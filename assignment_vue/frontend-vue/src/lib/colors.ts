export const classicColorKeys = ['black', 'blue', 'green', 'red', 'white'] as const

export type ClassicColorKey = typeof classicColorKeys[number]

const classicColorSet = new Set<string>(classicColorKeys)

export const classicColorSwatches = {
	black: '#111827',
	blue: '#5B6EE1',
	green: '#1FA34A',
	red: '#D92D48',
	white: '#FFFFFF',
} satisfies Record<ClassicColorKey, string>

export const normalizeColorName = (color: string) => color.trim().toLowerCase()

export const isClassicColorKey = (color: string): color is ClassicColorKey =>
	classicColorSet.has(normalizeColorName(color))

export const filterClassicColors = (colors: string[]): ClassicColorKey[] => {
	const seenColors = new Set<ClassicColorKey>()

	return colors.reduce<ClassicColorKey[]>((filteredColors, color) => {
		const normalizedColor = normalizeColorName(color)

		if (!isClassicColorKey(normalizedColor) || seenColors.has(normalizedColor)) {
			return filteredColors
		}

		seenColors.add(normalizedColor)
		filteredColors.push(normalizedColor)
		return filteredColors
	}, [])
}