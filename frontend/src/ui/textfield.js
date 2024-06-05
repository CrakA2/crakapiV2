
import { Button } from "./button"

export function Textfield(props) {
  return (
    <div className="flex items-center gap-2">
      <div className="bg-gray-100 dark:bg-gray-800 px-4 py-2 rounded-md flex-1 text-gray-900 dark:text-gray-100">
        {props.link}
      </div>
      <Button
        variant="ghost"
        size="icon"
        className="w-8 h-8 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md"
        onClick={() => {
          navigator.clipboard.writeText(props.link)
        }}
      >
        <CopyIcon className="w-4 h-4" />
        <span className="sr-only">Copy</span>
      </Button>
    </div>
  )
}

function CopyIcon(props) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <rect width="14" height="14" x="8" y="8" rx="2" ry="2" />
      <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
    </svg>
  )
}