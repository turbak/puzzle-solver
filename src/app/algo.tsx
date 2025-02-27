export const grid: string[][] = [
    ['jan', 'feb', 'mar', 'apr', 'may', 'jun'],
    ['jul', 'aug', 'sep', 'oct', 'nov', 'dec'],
    ['1', '2', '3', '4', '5', '6', '7'],
    ['8', '9', '10', '11', '12', '13','14'],
    ['15', '16', '17', '18', '19', '20', '21'],
    ['22', '23', '24', '25', '26', '27', '28'],
    ['29', '30', '31']
]

export type Month = 'jan' | 'feb' | 'mar' | 'apr' | 'may' | 'june' | 'july' | 'aug' | 'sep' | 'oct' | 'nov' | 'dec'

export type Piece = string[][]

const piece1: Piece = [
    ['{1}', ''],
    ['{1}', '{1}'],
    ['{1}', '{1}']
]

const piece2: Piece = [
    ['{2}', '{2}', '{2}', '{2}'],
    ['{2}', '',    '',    '']
]

const piece3: Piece = [
    ['{3}', '{3}', '{3}', ''],
    ['',    '',    '{3}', '{3}']
]

const piece4: Piece = [
    ['{4}', '{4}', '{4}'],
    ['{4}', '{4}', '{4}']
]

const piece5: Piece = [
    ['{5}', '{5}', ''],
    ['',    '{5}', ''],
    ['',    '{5}', '{5}'],
]

const piece6: Piece = [
    ['{6}', '{6}', '{6}', '{6}'],
    ['',    '{6}',  '',   '']
]

const piece7: Piece = [
    ['{7}', '{7}'],
    ['',    '{7}'],
    ['{7}', '{7}']
]

const piece8: Piece = [
    ['{8}',  '',    ''],
    ['{8}',  '',    ''],
    ['{8}',  '{8}', '{8}']
]

export const pieces = [piece1, piece2, piece3, piece4, piece5, piece6, piece7, piece8]

const DIRECTIONS = 4

const piecesRotations = pieces.map(piece => {
    const rotations = [piece]
    for (let i = 0; i < DIRECTIONS; i++) {
        rotations.push(rotateClockwisePiece(rotations[rotations.length - 1]))
    }

    return rotations
})

export function solve(month: Month, day: number): string[][] {
    if (day < 1 || day > 31) {
        throw new Error('Invalid day')
    }

    const initialGrid = copyGrid(grid)

    const piecesLeft = pieces.slice()
    const result = solveHelper(month, day.toString(), piecesLeft, initialGrid)
    if (!result) {
        throw new Error('No solution found')
    }

    return result
}

function canFindMonthAndDayOnGrid(month: Month, day: string, grid: string[][]): boolean {
    let dayFound = false
    let monthFound = false
    for (let i = 0; i < grid.length; i++) {
        for (let j = 0; j < grid[i].length; j++) {
            if (grid[i][j] == day) {
                dayFound = true
            } else if (grid[i][j] == month) {
                monthFound = true
            }

            if (dayFound && monthFound) {
                return true
            }
        }
    }

    return dayFound && monthFound
}

function solveHelper(month: Month, day: string, piecesLeft: Piece[], currentGrid: string[][]): string[][]|null {
    if (!canFindMonthAndDayOnGrid(month, day, currentGrid)) {
        return null
    }
    if (piecesLeft.length == 0) {
        return currentGrid
    }

    let currentPiece = piecesLeft[0]
    const newPiecesLeft = piecesLeft.slice(1)
    for(let turns = 0; turns < DIRECTIONS; turns++) {
        for (let rowStart = 0; rowStart < currentGrid.length; rowStart++) {
            for (let colStart = 0; colStart < currentGrid[rowStart].length; colStart++) {
                if (!canPlace(currentPiece, rowStart, colStart, currentGrid)) {
                    continue
                }


                placePiece(currentGrid, currentPiece, rowStart, colStart)

                const potenitalRes = solveHelper(month, day, newPiecesLeft, currentGrid)
                if (potenitalRes) {
                    return potenitalRes
                }

                unplacePiece(currentGrid, currentPiece, rowStart, colStart)
            }
        }

        currentPiece = rotateClockwisePiece(currentPiece)
    }

    return null
}

function copyGrid(grid: string[][]): string[][] {
    return grid.map(row => row.slice())
}

function placePiece(grid: string[][], piece: Piece, rowStart: number, colStart: number) {
    for (let i = 0; i < piece.length; i++) {
        for (let j = 0; j < piece[i].length; j++) {
            if (piece[i][j] == '') {
                continue
            }

            grid[rowStart+i][colStart+j] = piece[i][j]
        }
    }
}

function unplacePiece(currentGrid: string[][], piece: Piece, rowStart: number, colStart: number) {
    for (let i = 0; i < piece.length; i++) {
        for (let j = 0; j < piece[i].length; j++) {
            if (piece[i][j] == '') {
                continue
            }

            currentGrid[rowStart+i][colStart+j] = grid[rowStart+i][colStart+j]
        }
    }
}

function canPlace(piece: Piece, rowStart: number, colStart: number, grid: string[][]): boolean {
   for (let i = 0; i < piece.length; i++) {
        for (let j = 0; j < piece[i].length; j++) {
            if (piece[i][j] == '') {
                continue
            }

            if (grid.length - 1 < i + rowStart) {
                return false
            }

            if (grid[i+rowStart].length - 1 < j + colStart) {
                return false
            }
            
            if (grid[i+rowStart][j+colStart].includes('{')) {
                return false
            }
        }
    }

    return true
}

function rotateClockwisePiece(piece: Piece): Piece {
    const rotated = piece.map(row => row.slice())
    return rotated[0].map((_, index) => rotated.map(row => row[index]).reverse())
}

//console.table(solve('mar', 24))
