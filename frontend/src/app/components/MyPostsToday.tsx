import { Heart, MessageCircle, Clock } from "lucide-react";

const posts = [
  {
    id: 1,
    text: "Just launched a new feature for our product! Really excited to see how users respond to this update. Built with React and TypeScript 🚀",
    timestamp: "2 hours ago",
    likes: { count: 45, delta: 32 },
    replies: { count: 12, delta: 8 },
  },
  {
    id: 2,
    text: "Hot take: Dark mode should be the default for all developer tools. Who's with me?",
    timestamp: "4 hours ago",
    likes: { count: 78, delta: 54 },
    replies: { count: 23, delta: 15 },
  },
  {
    id: 3,
    text: "Working on some interesting glassmorphism designs today. The depth and layering effects are really coming together nicely.",
    timestamp: "6 hours ago",
    likes: { count: 34, delta: 21 },
    replies: { count: 5, delta: 4 },
  },
  {
    id: 4,
    text: "Quick tip: Always test your responsive designs on actual devices, not just browser dev tools. You'd be surprised at the differences!",
    timestamp: "8 hours ago",
    likes: { count: 56, delta: 38 },
    replies: { count: 9, delta: 7 },
  },
  {
    id: 5,
    text: "Coffee ☕ + Code 💻 = Perfect morning. What's your ideal coding setup?",
    timestamp: "10 hours ago",
    likes: { count: 89, delta: 67 },
    replies: { count: 18, delta: 12 },
  },
];

export default function MyPostsToday() {
  return (
    <div className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-6">
      <h2 className="text-xl font-bold text-white mb-4">My Posts Today</h2>
      <div className="space-y-4">
        {posts.map((post) => (
          <div
            key={post.id}
            className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-4 hover:bg-white/10 transition-all"
          >
            <p className="text-white mb-3 leading-relaxed">{post.text}</p>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-gray-400 text-sm">
                <Clock className="w-4 h-4" />
                <span>{post.timestamp}</span>
              </div>
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-1 px-2.5 py-1 rounded-full bg-pink-500/10 border border-pink-500/20">
                    <Heart className="w-3.5 h-3.5 text-pink-400" />
                    <span className="text-sm text-pink-400">{post.likes.count}</span>
                    <span className="text-xs text-pink-400/60">+{post.likes.delta}</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-1 px-2.5 py-1 rounded-full bg-blue-500/10 border border-blue-500/20">
                    <MessageCircle className="w-3.5 h-3.5 text-blue-400" />
                    <span className="text-sm text-blue-400">{post.replies.count}</span>
                    <span className="text-xs text-blue-400/60">+{post.replies.delta}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
