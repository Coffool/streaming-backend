import { User } from "lucide-react"

export default function UserDashboard() {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4 flex items-center">
        <User className="mr-2" /> User Dashboard
      </h1>
      <p>Welcome to your dashboard! Here you can manage your account and view your activities.</p>
    </div>
  )
}