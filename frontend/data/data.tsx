import {
  ArrowDown,
  ArrowRight,
  ArrowUp,
  CheckCircle,
  Circle,
  CircleAlert,
  CircleOff,
  HelpCircle,
  Timer,
} from "lucide-react"

export const labels = [
  {
    value: "bug",
    label: "Bug",
  },
  {
    value: "feature",
    label: "Feature",
  },
  // {
  //   value: "documentation",
  //   label: "Documentation",
  // },
]

export const statuses = [
  // {
  //   value: "backlog",
  //   label: "Backlog",
  //   icon: HelpCircle,
  // },
  // {
  //   value: "todo",
  //   label: "Todo",
  //   icon: Circle,
  // },
  // {
  //   value: "in progress",
  //   label: "In Progress",
  //   icon: Timer,
  // },
  // {
  //   value: "done",
  //   label: "Done",
  //   icon: CheckCircle,
  // },
  // {
  //   value: "canceled",
  //   label: "Canceled",
  //   icon: CircleOff,
  // },
  {
    value: "Pending",
    label: "Pending",
    icon: Circle,
  },
  {
    value: "In progress",
    label: "In Progress",
    icon: Timer,
  },
  {
    value: "Completed",
    label: "Completed",
    icon: CheckCircle,
  },
]

export const priorities = [
  // {
  //   label: "Low",
  //   value: "low",
  //   icon: ArrowDown,
  // },
  // {
  //   label: "Medium",
  //   value: "medium",
  //   icon: ArrowRight,
  // },
  // {
  //   label: "High",
  //   value: "high",
  //   icon: ArrowUp,
  // },
  // {
  //   label: "Critical",
  //   value: "critical",
  //   icon: CircleAlert,
  // },
  {
    label: "Low",
    value: "0.5",
    icon: ArrowDown,
  },
  {
    label: "Medium",
    value: "0.6",
    icon: ArrowRight,
  },
  {
    label: "High",
    value: "0.7",
    icon: ArrowUp,
  },
  {
    label: "Critical",
    value: "0.8",
    icon: CircleAlert,
  },
]
