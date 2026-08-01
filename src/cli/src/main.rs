use clap::{Parser, Subcommand};

/// 量潮学习云 CLI
#[derive(Parser)]
#[command(name = "qtcloud-learn", version, about)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// 打印版本信息
    Version,
}

fn main() {
    let cli = Cli::parse();
    match cli.command {
        Commands::Version => println!("qtcloud-learn {}", env!("CARGO_PKG_VERSION")),
    }
}
