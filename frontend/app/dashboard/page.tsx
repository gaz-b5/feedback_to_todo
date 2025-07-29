import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"


export default function Page() {
  return (
    <>
      <div className="w-full flex items-center justify-center h-full">
        <Alert variant="default" className="w-1/4">
          {/* <Terminal /> */}
          <AlertTitle>Heads up!</AlertTitle>
          <AlertDescription>
            Choose a project to see task-list 🙏.
          </AlertDescription>
        </Alert>
      </div>
    </>
  )
}
