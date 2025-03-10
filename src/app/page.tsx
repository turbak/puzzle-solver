"use client"

import { grid, solve, /*pieces, type Piece,*/ type Month } from './algo';
import { useEffect, useState } from 'react';
import { Analytics } from "@vercel/analytics/react"

export default function Home() {
    const [month, setMonth] = useState<Month>('jan')
    const [day, setDay] = useState(1)

    const [solGrid, setGrid] = useState<string[][] | null>(grid)

    const [displayFPS, setdisplayFPS] = useState(10)
    const [worker, setWorker] = useState<Worker | null>(null)

    const solvePuzzle = () => {
        fetch(`/api/solver?month=${month}&day=${day}`)
            .then((res) => res.json())
            .then((data) => setGrid(data))
            .then(() => console.log('Solved!'))
            .then(() => console.log(solGrid))
            .catch((err) => {
                console.error(err)

                const res = solve(month, day)
                if (res === null) {
                    console.error('No solution found')
                    setGrid(null)
                    return
                }
                setGrid(res)
            });
    }

    useEffect(() => {
        return () => {
            if (worker) {
                worker.terminate();
            }
            setGrid(null);
        };
    }, [worker]);


    const solvePuzzleWithDisplay = () => {
        if (typeof Worker == 'undefined') {
            console.error('Web worker not supported')
            return
        }

        const worker = new Worker(new URL('./solverWorker.ts', import.meta.url));
        setWorker(worker);

        worker.onmessage = (event) => {
            if (event.data.error) {
                console.error(event.data.error);
            } else if (event.data.grid) {
                setGrid(event.data.grid);
            } else if (event.data.result) {
                setGrid(event.data.result);
                if (worker) {
                    worker.terminate();
                    setWorker(null);
                }
            }
        };

        worker.postMessage({
            month: month,
            day: day,
            throttleInterval: 1000 / displayFPS,
        });
    }

    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-6">
            <h1 className="text-4xl font-bold text-gray-800 mb-6">Puzzle Solver</h1>
            <div className="bg-white shadow-lg rounded-2xl p-6 w-full max-w-lg flex flex-col items-center space-y-4">

                <GridView grid={solGrid || grid} />

                <div className="flex flex-col space-y-2 w-full">
                    <label className="flex flex-col text-gray-700 font-semibold">
                        Month:
                        <select
                            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mt-1"
                            value={month}
                            onChange={e => setMonth(e.target.value as Month)}
                        >
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
                    </label>

                    <label className="flex flex-col text-gray-700 font-semibold">
                        Day:
                        <input
                            className="bg-white text-gray-900 font-bold py-2 px-4 rounded border border-gray-300 mt-1"
                            type="number"
                            value={day}
                            pattern="^([1-9])|([1-3][0-9])$"
                            onChange={e => setDay(parseInt(e.target.value))}
                        />
                    </label>
                </div>
                <div className='flex gap-4 justify-center'>
                    <button
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-xl shadow-md transition-all"
                        onClick={() => {
                            if (worker) {
                                worker.terminate()
                                setWorker(null)
                            }
                            solvePuzzle()
                        }}
                    >
                        Solve
                    </button>

                    <button
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-xl shadow-md transition-all"
                        onClick={() => {
                            if (worker) {
                                worker.terminate()
                                setWorker(null)
                            }
                            setGrid(null)
                        }}
                    >
                        Clear
                    </button>

                    <button
                        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-xl shadow-md transition-all"
                        onClick={() => {
                            if (worker) {
                                worker.terminate()
                                setWorker(null)
                                setGrid(null)
                            }
                            solvePuzzleWithDisplay()
                        }}
                    >
                        Solve with display
                    </button>
                </div>
            </div>
            <label className="flex flex-col text-gray-700 font-semibold">
                Display Frames Per Second:
            </label>
            <input
                className="bg-white text-gray-900 font-bold py-2 px-4 rounded border border-gray-300"
                type="number"
                value={displayFPS}
                onChange={e => {
                    setdisplayFPS(parseInt(e.target.value))
                    if (worker) {
                        worker.terminate()
                        setWorker(null)
                        setGrid(null)
                        solvePuzzleWithDisplay()
                    }
                }}
            />
            <Analytics />
        </div>
    )
}

/*
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
            */

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

/*
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
            */

function GridView(props: {
    grid: string[][],
}) {
    return (
        <div className='grid grid-cols-7 gap-1 rounded-lg shadow-lg w-fit border-1 border-gray-400 p-1'>
            {Array.from({ length: 49 }).map((_, index) => {
                const i = Math.floor(index / 7);
                const j = index % 7;
                const cell = props.grid[i]?.[j];

                return (
                    <div
                        key={`${i}-${j}`}
                        className={`w-10 h-10 flex items-center justify-center font-bold mt-1 rounded-lg 
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
