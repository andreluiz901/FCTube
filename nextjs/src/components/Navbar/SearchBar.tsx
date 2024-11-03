"use client";

import { useSearchParams } from "next/navigation";

export function SearchBar() {
  const searchParams = useSearchParams();
  return (
    <form action="/" method="get">
      <input
        type="text"
        id="search"
        name="search"
        placeholder="Pesquisar"
        className="w-full py-2 px-4 bg-primary border-secondary border-2 rounded-full focus:outline-none focus:ring-2 focus:ring-secondary text-primary"
        defaultValue={searchParams.get("search") || ""}
      />
      <button className="absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-600">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          className="size-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z"
          />
        </svg>
      </button>
    </form>
  );
}
