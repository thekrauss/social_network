// Sidebar Component
"use client";
import React, { useState } from "react";
import { IconMenu2, IconArrowLeft, IconBrandTabler, IconSettings, IconUserBolt, IconHelpCircle, IconLogout, IconBell, IconMessageCircle } from "@tabler/icons-react";
import Link from "next/link";
import { useRouter } from "next/router";

export function SidebarDemo() {
  const links = [
    {
      label: "Dashboard",
      href: "/admin",
      icon: <IconBrandTabler className="text-gray-400 h-6 w-6" />,
    },
    {
      label: "Profile",
      href: "/profile",
      icon: <IconUserBolt className="text-gray-400 h-6 w-6" />,
    },
    {
      label: "Settings",
      href: "/settings",
      icon: <IconSettings className="text-gray-400 h-6 w-6" />,
    },
  ];

  const [open, setOpen] = useState(false);
  const router = useRouter();

  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:8079/logout", {
        method: "POST",
        credentials: "include",
      });
      if (response.ok) {
        router.push("/login");
      } else {
        console.error("Erreur lors de la déconnexion");
      }
    } catch (error) {
      console.error("Erreur de connexion au backend:", error);
    }
  };

  return (
    <div className="flex">
      {/* Sidebar */}
      <div className={`fixed top-0 left-0 h-full bg-gray-900 text-white shadow-xl z-50 lg:relative lg:w-64 transition-transform duration-300 ${open ? "translate-x-0" : "-translate-x-full"} lg:translate-x-0`}>
        <Sidebar open={open} setOpen={setOpen}>
          <SidebarBody className="flex flex-col justify-between gap-6">
            <Logo open={open} setOpen={setOpen} />
            <div className="mt-8 flex flex-col gap-6">
              {links.map((link, idx) => (
                <SidebarLink key={idx} link={link} open={open} />
              ))}
            </div>
            <div className="mt-auto flex flex-col gap-2 border-t border-gray-700 pt-4">
              <SidebarLink link={{ label: "Aide", href: "/help", icon: <IconHelpCircle className="text-gray-400 h-6 w-6" /> }} open={open} />
              <SidebarLink link={{ label: "Support", href: "/support", icon: <IconHelpCircle className="text-gray-400 h-6 w-6" /> }} open={open} />
              <button onClick={handleLogout} className="flex items-center space-x-3 p-3 hover:bg-red-600 rounded-lg transition-all duration-300 hover:text-white text-gray-300">
                <IconLogout className="h-6 w-6 text-gray-400" />
                <span className="text-sm">Déconnexion</span>
              </button>
            </div>
          </SidebarBody>
        </Sidebar>
      </div>

      {/* Sidebar Toggle Button for mobile */}
      <div className="p-4 lg:hidden">
        <button onClick={() => setOpen(!open)} className="text-white bg-gray-800 p-2 rounded-lg shadow-lg hover:bg-cyan-600 transition-all duration-300">
          {open ? <IconArrowLeft className="h-6 w-6" /> : <IconMenu2 className="h-6 w-6" />}
        </button>
      </div>
    </div>
  );
}

export const Logo = ({ open, setOpen }) => (
  <button onClick={() => setOpen(!open)} className="flex items-center space-x-2 text-sm text-white py-6 border-b border-gray-700 w-full justify-center">
    <div className="h-8 w-8 bg-cyan-500 rounded-lg flex-shrink-0 shadow-lg" />
    <span className="text-lg font-semibold">MyApp</span>
  </button>
);

export function Sidebar({ open, setOpen, children }) {
  return <aside>{children}</aside>;
}

export function SidebarBody({ children, className }) {
  return <div className={`px-4 flex flex-col h-full ${className}`}>{children}</div>;
}

export function SidebarLink({ link }) {
  return (
    <Link href={link.href} className="sidebar-link flex items-center space-x-3 p-3 hover:bg-cyan-600 rounded-lg transition-all duration-300 hover:text-white">
      {link.icon}
      <span className="text-sm text-gray-300">{link.label}</span>
    </Link>
  );
}