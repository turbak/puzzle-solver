"use client"

import { solve, grid, pieces, type Piece, type Month } from './algo';
import { useState } from 'react';

export default function Home() {
    const [month, setMonth] = useState<Month>('jan')
    const [day, setDay] = useState(1)

    const [solution, setSolution] = useState<string[][] | null>(grid)

    const solvePuzzle = () => {
        setSolution(grid)
        try {
            setSolution(solve(month, day))
        } catch (e) {
            alert(e)
        }
    }

    return <>
        <h1>Puzzle Solver</h1>
        <div className=''>
            <PieceSelector pieces={pieces} setPieces={console.log} />
            <GridView grid={solution || grid} />
            <label>
                Month:
                <select value={month} onChange={e => setMonth(e.target.value as Month)}>
                    <option value="jan">January</option>
                    <option value="feb">February</option>
                    <option value="mar">March</option>
                    <option value="apr">April</option>
                    <option value="may">May</option>
                    <option value="jun">June</option>
                    <option value="jul">July</option>
                    <option value="aug">August</option>
                    <option value="sep">September</option>
                    <option value="oct">October</option>
                    <option value="nov">November</option>
                    <option value="dec">December</option>
                </select>
                Day:
                <input type="number" value={day} pattern='^([1-9])|([1-3][0-9])$' onChange={e => setDay(parseInt(e.target.value))} />
            </label>
            <button onClick={solvePuzzle}>Solve</button>
        </div>
    </>
}


function PieceView(props: {
    piece: Piece,
    id: number,
}) {
    return (
        <div className="grid grid-cols-4 gap-1">
            {props.piece.map((row, i) => (
                <div key={i} className="flex gap-1">
                    {row.map((cell, j) => (
                        <div key={j} className={`${colorFromId(props.id)} w-8 h-8`} />
                    ))}
                </div>
            ))}
        </div>
    )
}

function colorFromId(id: number): string {
    switch (id) {
        case 0: return 'bg-red-500'
        case 1: return 'bg-yellow-500'
        case 2: return 'bg-green-500'
        case 3: return 'bg-blue-500'
        case 4: return 'bg-purple-500'
        case 5: return 'bg-pink-500'
        case 6: return 'bg-indigo-500'
        case 7: return 'bg-cyan-500'
        case 8: return 'bg-rose-500'
        default: return 'bg-gray-500'
    }
}


function PieceSelector(props: {
    pieces: Piece[],
    setPieces: (pieces: Piece[]) => void,
}) {
    return (
        <div className="grid grid-cols-4 gap-4">
            {props.pieces.map((piece, i) => (
                <PieceView key={i} piece={piece} id={i} />
            ))}
        </div>
    )
}

function GridView(props: {
    grid: string[][],
}) {
    return (
        <div className='grid grid-cols-7 gap-1 w-fit border-2 border-gray-400 p-1'>
            {Array.from({ length: 49 }).map((_, index) => {
                const i = Math.floor(index / 7);
                const j = index % 7;
                const cell = props.grid[i]?.[j];

                return (
                    <div
                        key={`${i}-${j}`}
                        className={`w-8 h-8 flex items-center justify-center
                            ${cell?.includes('{')
                                ? colorFromId(parseInt(cell.replace('{', '').replace('}', '')))
                                : 'bg-gray-500'
                            }`}
                    >
                        {cell?.includes('{') ? '' : cell}
                    </div>
                );
            })}
        </div>
    );
}
