import { Twitter, Home, BarChart3, Settings, LogOut, User } from "lucide-react";
import { useNavigate } from "react-router";

export default function Sidebar() {
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate("/");
  };

  return (
    <>
      {/* Desktop sidebar */}
      <aside className="hidden lg:flex flex-col w-64 backdrop-blur-2xl bg-white/5 border-r border-white/10 p-6">
        {/* Logo */}
        <div className="flex items-center gap-3 mb-8">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center">
            <Twitter className="w-5 h-5 text-white" />
          </div>
          <div>
            <h2 className="font-bold text-white">X Daily Digest</h2>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 space-y-2">
          <a
            href="#"
            className="flex items-center gap-3 px-4 py-3 rounded-[16px] bg-white/10 text-white"
          >
            <Home className="w-5 h-5" />
            <span>Dashboard</span>
          </a>
          <a
            href="#"
            className="flex items-center gap-3 px-4 py-3 rounded-[16px] text-gray-400 hover:bg-white/5 hover:text-white transition-all"
          >
            <BarChart3 className="w-5 h-5" />
            <span>Analytics</span>
          </a>
          <a
            href="#"
            className="flex items-center gap-3 px-4 py-3 rounded-[16px] text-gray-400 hover:bg-white/5 hover:text-white transition-all"
          >
            <User className="w-5 h-5" />
            <span>Profile</span>
          </a>
          <a
            href="#"
            className="flex items-center gap-3 px-4 py-3 rounded-[16px] text-gray-400 hover:bg-white/5 hover:text-white transition-all"
          >
            <Settings className="w-5 h-5" />
            <span>Settings</span>
          </a>
        </nav>

        {/* User profile & logout */}
        <div className="pt-4 border-t border-white/10 space-y-2">
          <div className="flex items-center gap-3 px-4 py-3">
            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-500" />
            <div className="flex-1 min-w-0">
              <p className="text-white text-sm font-medium truncate">@yourhandle</p>
              <p className="text-gray-400 text-xs truncate">Your Name</p>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-4 py-3 rounded-[16px] text-gray-400 hover:bg-white/5 hover:text-red-400 transition-all w-full"
          >
            <LogOut className="w-5 h-5" />
            <span>Logout</span>
          </button>
        </div>
      </aside>

      {/* Mobile bottom navigation */}
      <div className="lg:hidden fixed bottom-0 left-0 right-0 backdrop-blur-2xl bg-white/5 border-t border-white/10 z-50">
        <nav className="flex items-center justify-around p-4">
          <a href="#" className="flex flex-col items-center gap-1 text-white">
            <Home className="w-6 h-6" />
            <span className="text-xs">Home</span>
          </a>
          <a href="#" className="flex flex-col items-center gap-1 text-gray-400">
            <BarChart3 className="w-6 h-6" />
            <span className="text-xs">Analytics</span>
          </a>
          <a href="#" className="flex flex-col items-center gap-1 text-gray-400">
            <User className="w-6 h-6" />
            <span className="text-xs">Profile</span>
          </a>
          <a href="#" className="flex flex-col items-center gap-1 text-gray-400">
            <Settings className="w-6 h-6" />
            <span className="text-xs">Settings</span>
          </a>
        </nav>
      </div>
    </>
  );
}
