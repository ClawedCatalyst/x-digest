import { Twitter, Shield } from "lucide-react";

export default function LoginPage() {
  const handleLogin = () => {
    const backendBaseUrl =
      import.meta.env.VITE_BACKEND_URL ?? "https://699e-2405-201-6007-80b6-d8d9-b4dd-6e34-6512.ngrok-free.app";
    window.location.href = `${backendBaseUrl}/auth/x/start`;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0a0a0f] via-[#141420] to-[#1a1a2e] flex items-center justify-center p-4">
      {/* Background gradient orbs */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl" />
      </div>

      <div className="relative w-full max-w-md">
        {/* Glassmorphism card */}
        <div className="backdrop-blur-2xl bg-white/5 border border-white/10 rounded-[16px] p-8 shadow-2xl">
          {/* Logo and branding */}
          <div className="flex items-center justify-center mb-8">
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center">
                <Twitter className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold text-white">X Daily Digest</h1>
              </div>
            </div>
          </div>

          {/* Tagline */}
          <p className="text-center text-gray-300 mb-8 leading-relaxed">
            Your daily summary of what happened on your posts
          </p>

          {/* Continue with X button */}
          <button
            onClick={handleLogin}
            className="w-full bg-white text-black rounded-[16px] px-6 py-4 flex items-center justify-center gap-3 hover:bg-gray-100 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-[1.02]"
          >
            <Twitter className="w-5 h-5" />
            <span className="font-semibold">Continue with X</span>
          </button>

          {/* Privacy note */}
          <div className="mt-6 flex items-start gap-2 text-sm text-gray-400">
            <Shield className="w-4 h-4 mt-0.5 flex-shrink-0 text-green-400" />
            <p className="leading-relaxed">
              We only read your public activity and engagement on your posts. Your privacy is important to us.
            </p>
          </div>
        </div>

        {/* Footer text */}
        <p className="text-center text-gray-500 text-sm mt-6">
          Stay updated with your daily X activity insights
        </p>
      </div>
    </div>
  );
}
