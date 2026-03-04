import { MessageCircle, ExternalLink, Clock } from "lucide-react";

const replies = [
  {
    id: 1,
    author: "@johndoe",
    authorAvatar: 1,
    replyText: "This is amazing! I've been waiting for a feature like this for months. Can't wait to try it out!",
    time: "1h ago",
    postPreview: "Just launched a new feature...",
  },
  {
    id: 2,
    author: "@sarahsmith",
    authorAvatar: 2,
    replyText: "100% agreed! Dark mode is so much easier on the eyes, especially during late night coding sessions.",
    time: "3h ago",
    postPreview: "Hot take: Dark mode should be...",
  },
  {
    id: 3,
    author: "@mikejones",
    authorAvatar: 3,
    replyText: "Love the depth effects! What tools are you using for this? Figma?",
    time: "5h ago",
    postPreview: "Working on some interesting...",
  },
  {
    id: 4,
    author: "@emilybrown",
    authorAvatar: 4,
    replyText: "So true! I found a critical bug on mobile that didn't show up in Chrome DevTools at all.",
    time: "7h ago",
    postPreview: "Quick tip: Always test your...",
  },
  {
    id: 5,
    author: "@alexwilson",
    authorAvatar: 5,
    replyText: "Mechanical keyboard + standing desk + good headphones = chef's kiss 👨‍🍳",
    time: "9h ago",
    postPreview: "Coffee + Code = Perfect morning...",
  },
];

const getAvatarColor = (id: number) => {
  const colors = [
    "from-blue-500 to-cyan-500",
    "from-purple-500 to-pink-500",
    "from-orange-500 to-red-500",
    "from-green-500 to-emerald-500",
    "from-indigo-500 to-purple-500",
  ];
  return colors[id % colors.length];
};

export default function NewRepliesToday() {
  return (
    <div className="backdrop-blur-xl bg-white/5 border border-white/10 rounded-[16px] p-6 h-fit">
      <div className="flex items-center gap-2 mb-4">
        <MessageCircle className="w-5 h-5 text-blue-400" />
        <h2 className="text-xl font-bold text-white">New Replies Today</h2>
      </div>
      <div className="space-y-4">
        {replies.map((reply) => (
          <div key={reply.id} className="pb-4 border-b border-white/10 last:border-0 last:pb-0">
            <div className="flex items-start gap-3 mb-2">
              <div
                className={`w-10 h-10 rounded-full bg-gradient-to-br ${getAvatarColor(
                  reply.authorAvatar
                )} flex-shrink-0`}
              />
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-white font-medium">{reply.author}</span>
                  <span className="text-xs text-gray-500">·</span>
                  <div className="flex items-center gap-1 text-gray-400 text-xs">
                    <Clock className="w-3 h-3" />
                    <span>{reply.time}</span>
                  </div>
                </div>
                <p className="text-sm text-gray-300 leading-relaxed mb-2">{reply.replyText}</p>
                <p className="text-xs text-gray-500 mb-2 line-clamp-1">
                  Reply to: {reply.postPreview}
                </p>
                <a
                  href="#"
                  className="inline-flex items-center gap-1 text-xs text-blue-400 hover:text-blue-300 transition-colors"
                >
                  <ExternalLink className="w-3 h-3" />
                  <span>View thread</span>
                </a>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
