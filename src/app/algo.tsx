export const grid: string[][] = [
    ['jan', 'feb', 'mar', 'apr', 'may', 'jun'],
    ['jul', 'aug', 'sep', 'oct', 'nov', 'dec'],
    ['1', '2', '3', '4', '5', '6', '7'],
    ['8', '9', '10', '11', '12', '13', '14'],
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
    ['{2}', '', '', '']
]

const piece3: Piece = [
    ['{3}', '{3}', '{3}', ''],
    ['', '', '{3}', '{3}']
]

const piece4: Piece = [
    ['{4}', '{4}', '{4}'],
    ['{4}', '{4}', '{4}']
]

const piece5: Piece = [
    ['{5}', '{5}', ''],
    ['', '{5}', ''],
    ['', '{5}', '{5}'],
]

const piece6: Piece = [
    ['{6}', '{6}', '{6}', '{6}'],
    ['', '{6}', '', '']
]

const piece7: Piece = [
    ['{7}', '{7}'],
    ['', '{7}'],
    ['{7}', '{7}']
]

const piece8: Piece = [
    ['{8}', '', ''],
    ['{8}', '', ''],
    ['{8}', '{8}', '{8}']
]

export const pieces = [piece1, piece2, piece3, piece4, piece5, piece6, piece7, piece8]
    .sort((a, b) => {
        const sizeA = a.flat().filter(c => c !== '').length;
        const sizeB = b.flat().filter(c => c !== '').length;
        return sizeB - sizeA;
    });

const MAX_ROTATIONS = 4

const pieceRotationsCache = new Map<Piece, Piece[]>();

function getRotations(piece: Piece): Piece[] {
    if (pieceRotationsCache.has(piece)) {
        return pieceRotationsCache.get(piece)!;
    }

    const seen = new Set<string>();
    const rotations: Piece[] = [];
    let current = piece;

    for (let i = 0; i < MAX_ROTATIONS; i++) {
        const variants = [
            current,
            flipPiece(current)
        ];

        for (const variant of variants) {
            const key = JSON.stringify(variant);
            if (!seen.has(key)) {
                seen.add(key);
                rotations.push(variant);
            }
        }

        const next = rotateClockwisePiece(current);
        if (JSON.stringify(next) === JSON.stringify(piece)) break;
        current = next;
    }

    pieceRotationsCache.set(piece, rotations);
    return rotations;
}

export function solve(month: Month, day: number): string[][] {
    if (day < 1 || day > 31) {
        throw new Error('Invalid day')
    }

    const initialGrid = copyGrid(grid)

    const piecesLeft = pieces.slice()
    const result = solveHelper(month, day.toString(), 0, initialGrid)
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

function solveHelper(month: Month, day: string, pieceIdx: number, solutionGrid: string[][]): string[][] | null {
    if (!canFindMonthAndDayOnGrid(month, day, solutionGrid)) {
        return null
    }
    if (pieceIdx >= pieces.length) {
        return solutionGrid
    }

    for (const rotatedPiece of getRotations(pieces[pieceIdx])) {
        for (let rowStart = 0; rowStart < solutionGrid.length; rowStart++) {
            for (let colStart = 0; colStart < solutionGrid[rowStart].length; colStart++) {
                if (!canPlace(rotatedPiece, rowStart, colStart, solutionGrid)) {
                    continue
                }

                placePiece(solutionGrid, rotatedPiece, rowStart, colStart)

                const potenitalRes = solveHelper(month, day, pieceIdx + 1, solutionGrid)
                if (potenitalRes) {
                    return potenitalRes
                }

                unplacePiece(solutionGrid, rotatedPiece, rowStart, colStart)
            }
        }
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

            const row = i + rowStart
            const col = j + colStart

            grid[row][col] = piece[i][j]
        }
    }
}

function unplacePiece(currentGrid: string[][], piece: Piece, rowStart: number, colStart: number) {
    for (let i = 0; i < piece.length; i++) {
        for (let j = 0; j < piece[i].length; j++) {
            if (piece[i][j] == '') {
                continue
            }

            currentGrid[rowStart + i][colStart + j] = grid[rowStart + i][colStart + j]
        }
    }
}

function canPlace(piece: Piece, rowStart: number, colStart: number, grid: string[][]): boolean {
    for (let i = 0; i < piece.length; i++) {
        for (let j = 0; j < piece[i].length; j++) {
            if (piece[i][j] == '') {
                continue
            }

            const row = i + rowStart
            const col = j + colStart

            if (grid.length - 1 < row) {
                return false
            }

            if (grid[row].length - 1 < col) {
                return false
            }

            if (grid[row][col].includes('{')) {
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

function flipPiece(piece: Piece): Piece {
    return piece.map(row => row.slice()).reverse()
}

console.table(solve('feb', 26))
