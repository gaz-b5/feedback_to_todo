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
    value: "In Progress",
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
    value: "1",
    icon: ArrowDown,
  },
  {
    label: "Medium",
    value: "2",
    icon: ArrowRight,
  },
  {
    label: "High",
    value: "3",
    icon: ArrowUp,
  },
  {
    label: "Critical",
    value: "4",
    icon: CircleAlert,
  },
]

export const roles = [
  {
    value: "admin",
    label: "Admin",
  },
  {
    value: "member",
    label: "Member",
  },
]
